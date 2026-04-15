package handlers

import (
	"net/http"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListEnrollments(c *gin.Context) {
	userID, _ := c.Get("userId")
	role, _ := c.Get("userRole")
	userRole := role.(models.UserRole)

	query := db.DB.Preload("Module.Category").Preload("User")

	if userRole == models.RoleLearner {
		query = query.Where("user_id = ?", userID)
	} else if uid := c.Query("userId"); uid != "" {
		query = query.Where("user_id = ?", uid)
	}

	if moduleID := c.Query("moduleId"); moduleID != "" {
		query = query.Where("module_id = ?", moduleID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	var enrollments []models.Enrollment
	query.Order("enrolled_at DESC").Find(&enrollments)
	c.JSON(http.StatusOK, enrollments)
}

func CreateEnrollment(c *gin.Context) {
	var req struct {
		UserID   string  `json:"userId" binding:"required"`
		ModuleID string  `json:"moduleId" binding:"required"`
		DueDate  *string `json:"dueDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignedBy, _ := c.Get("userId")

	enrollment := models.Enrollment{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		ModuleID:   req.ModuleID,
		Status:     models.EnrollmentNotStarted,
		AssignedBy: ptrString(assignedBy.(string)),
	}

	if err := db.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Enrollment already exists"})
		return
	}

	db.DB.Preload("Module").Preload("User").First(&enrollment, "id = ?", enrollment.ID)
	c.JSON(http.StatusCreated, enrollment)
}

func BulkEnroll(c *gin.Context) {
	var req struct {
		UserIDs  []string `json:"userIds" binding:"required"`
		ModuleID string   `json:"moduleId" binding:"required"`
		DueDate  *string  `json:"dueDate"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	assignedBy, _ := c.Get("userId")
	created := 0
	skipped := 0

	for _, uid := range req.UserIDs {
		enrollment := models.Enrollment{
			ID:         uuid.New().String(),
			UserID:     uid,
			ModuleID:   req.ModuleID,
			Status:     models.EnrollmentNotStarted,
			AssignedBy: ptrString(assignedBy.(string)),
		}
		if err := db.DB.Create(&enrollment).Error; err != nil {
			skipped++
		} else {
			created++
		}
	}

	c.JSON(http.StatusOK, gin.H{"created": created, "skipped": skipped})
}

func UpdateEnrollment(c *gin.Context) {
	id := c.Param("id")
	var enrollment models.Enrollment
	if err := db.DB.First(&enrollment, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Enrollment not found"})
		return
	}

	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	db.DB.Model(&enrollment).Updates(req)
	c.JSON(http.StatusOK, enrollment)
}

func DeleteEnrollment(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Enrollment{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "Enrollment deleted"})
}

func SelfEnroll(c *gin.Context) {
	userID, _ := c.Get("userId")
	var req struct {
		ModuleID string `json:"moduleId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify module is published
	var module models.Module
	if err := db.DB.Where("id = ? AND status = ?", req.ModuleID, models.StatusPublished).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found or not published"})
		return
	}

	enrollment := models.Enrollment{
		ID:       uuid.New().String(),
		UserID:   userID.(string),
		ModuleID: req.ModuleID,
		Status:   models.EnrollmentNotStarted,
	}

	if err := db.DB.Create(&enrollment).Error; err != nil {
		// Already enrolled - just return existing
		var existing models.Enrollment
		db.DB.Where("user_id = ? AND module_id = ?", userID, req.ModuleID).
			Preload("Module", func(tx *gorm.DB) *gorm.DB { return tx }).
			First(&existing)
		c.JSON(http.StatusOK, existing)
		return
	}

	db.DB.Preload("Module").First(&enrollment, "id = ?", enrollment.ID)
	c.JSON(http.StatusCreated, enrollment)
}

func ptrString(s string) *string {
	return &s
}
