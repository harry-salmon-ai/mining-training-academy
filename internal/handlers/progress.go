package handlers

import (
	"math"
	"net/http"
	"time"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RecordSlideProgress(c *gin.Context) {
	uid, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	var req struct {
		SlideID   string `json:"slideId" binding:"required"`
		TimeSpent int    `json:"timeSpent"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()

	var existing models.SlideProgress
	err := db.DB.Where("user_id = ? AND slide_id = ?", uid, req.SlideID).First(&existing).Error

	if err == nil {
		db.DB.Model(&existing).Updates(map[string]interface{}{
			"visit_count": existing.VisitCount + 1,
			"last_visit":  now,
			"completed":   true,
			"time_spent":  existing.TimeSpent + req.TimeSpent,
		})
	} else {
		db.DB.Create(&models.SlideProgress{
			ID:        uuid.New().String(),
			UserID:    uid,
			SlideID:   req.SlideID,
			Completed: true,
			TimeSpent: req.TimeSpent,
		})
	}

	recalculateModuleProgress(uid, req.SlideID)

	c.JSON(http.StatusOK, gin.H{"message": "Progress recorded"})
}

func GetModuleProgress(c *gin.Context) {
	moduleID := c.Param("id")
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var enrollment models.Enrollment
	err := db.DB.Where("user_id = ? AND module_id = ?", userID, moduleID).
		Preload("Module.Category").
		First(&enrollment).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Not enrolled"})
		return
	}

	var completedSlides int64
	db.DB.Model(&models.SlideProgress{}).
		Joins("JOIN slides ON slides.id = slide_progresses.slide_id").
		Joins("JOIN sections ON sections.id = slides.section_id").
		Where("sections.module_id = ? AND slide_progresses.user_id = ? AND slide_progresses.completed = true", moduleID, userID).
		Count(&completedSlides)

	var totalSlides int64
	db.DB.Model(&models.Slide{}).
		Joins("JOIN sections ON sections.id = slides.section_id").
		Where("sections.module_id = ?", moduleID).
		Count(&totalSlides)

	c.JSON(http.StatusOK, gin.H{
		"enrollment":      enrollmentToListJSON(&enrollment, false),
		"totalSlides":     totalSlides,
		"completedSlides": completedSlides,
		"progress":        enrollment.Progress,
	})
}

func GetProgressOverview(c *gin.Context) {
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var enrollments []models.Enrollment
	db.DB.Where("user_id = ?", userID).Preload("Module.Category").Find(&enrollments)

	completed := 0
	inProgress := 0
	for _, e := range enrollments {
		if e.Status == models.EnrollmentCompleted {
			completed++
		} else if e.Status == models.EnrollmentInProgress {
			inProgress++
		}
	}

	var certCount int64
	db.DB.Model(&models.Certificate{}).Where("user_id = ?", userID).Count(&certCount)

	details := make([]gin.H, 0, len(enrollments))
	for i := range enrollments {
		details = append(details, enrollmentToListJSON(&enrollments[i], false))
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments":      len(enrollments),
		"completedModules": completed,
		"inProgress":       inProgress,
		"certificates":     certCount,
		"details":          details,
	})
}

func recalculateModuleProgress(userID, slideID string) {
	var slide models.Slide
	if err := db.DB.Preload("Section").First(&slide, "id = ?", slideID).Error; err != nil {
		return
	}

	moduleID := slide.Section.ModuleID

	var totalSlides int64
	db.DB.Model(&models.Slide{}).
		Joins("JOIN sections ON sections.id = slides.section_id").
		Where("sections.module_id = ?", moduleID).
		Count(&totalSlides)

	var completedSlides int64
	db.DB.Model(&models.SlideProgress{}).
		Joins("JOIN slides ON slides.id = slide_progresses.slide_id").
		Joins("JOIN sections ON sections.id = slides.section_id").
		Where("sections.module_id = ? AND slide_progresses.user_id = ? AND slide_progresses.completed = true", moduleID, userID).
		Count(&completedSlides)

	progress := float64(0)
	if totalSlides > 0 {
		progress = math.Round(float64(completedSlides) / float64(totalSlides) * 100)
	}

	updates := map[string]interface{}{
		"progress": progress,
	}
	if progress > 0 {
		now := time.Now()
		updates["started_at"] = now
		updates["status"] = models.EnrollmentInProgress
	}
	if progress == 100 {
		now := time.Now()
		updates["completed_at"] = now
		updates["status"] = models.EnrollmentCompleted
	}

	db.DB.Model(&models.Enrollment{}).
		Where("user_id = ? AND module_id = ?", userID, moduleID).
		Updates(updates)
}
