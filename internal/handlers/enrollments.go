package handlers

import (
	"net/http"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func enrollmentToListJSON(e *models.Enrollment, includeUser bool) gin.H {
	h := gin.H{
		"id":          e.ID,
		"userId":      e.UserID,
		"moduleId":    e.ModuleID,
		"status":      e.Status,
		"progress":    e.Progress,
		"enrolledAt":  e.EnrolledAt,
		"startedAt":   e.StartedAt,
		"completedAt": e.CompletedAt,
		"dueDate":     e.DueDate,
	}
	mod := gin.H{
		"id":    e.Module.ID,
		"title": e.Module.Title,
		"slug":  e.Module.Slug,
	}
	if e.Module.CategoryID != "" {
		mod["category"] = gin.H{
			"id":    e.Module.Category.ID,
			"name":  e.Module.Category.Name,
			"slug":  e.Module.Category.Slug,
			"color": e.Module.Category.Color,
		}
	}
	h["module"] = mod
	if includeUser {
		h["user"] = newAuthUserResponse(&e.User)
	}
	return h
}

func ListEnrollments(c *gin.Context) {
	userID, uidOK := contextUserID(c)
	userRole, roleOK := contextUserRole(c)
	if !uidOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if !roleOK {
		userRole = models.RoleLearner
	}

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
	includeUser := userRole != models.RoleLearner
	out := make([]gin.H, 0, len(enrollments))
	for i := range enrollments {
		out = append(out, enrollmentToListJSON(&enrollments[i], includeUser))
	}
	c.JSON(http.StatusOK, out)
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

	assignedBy, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	enrollment := models.Enrollment{
		ID:         uuid.New().String(),
		UserID:     req.UserID,
		ModuleID:   req.ModuleID,
		Status:     models.EnrollmentNotStarted,
		AssignedBy: ptrString(assignedBy),
	}

	if err := db.DB.Create(&enrollment).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Enrollment already exists"})
		return
	}

	db.DB.Preload("Module.Category").Preload("User").First(&enrollment, "id = ?", enrollment.ID)
	c.JSON(http.StatusCreated, enrollmentToListJSON(&enrollment, true))
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

	assignedBy, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	created := 0
	skipped := 0

	for _, uid := range req.UserIDs {
		enrollment := models.Enrollment{
			ID:         uuid.New().String(),
			UserID:     uid,
			ModuleID:   req.ModuleID,
			Status:     models.EnrollmentNotStarted,
			AssignedBy: ptrString(assignedBy),
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
	db.DB.Preload("Module.Category").Preload("User").First(&enrollment, "id = ?", id)
	ur, roleOK := contextUserRole(c)
	includeUser := roleOK && ur != models.RoleLearner
	c.JSON(http.StatusOK, enrollmentToListJSON(&enrollment, includeUser))
}

func DeleteEnrollment(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Enrollment{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "Enrollment deleted"})
}

func SelfEnroll(c *gin.Context) {
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
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
		UserID:   userID,
		ModuleID: req.ModuleID,
		Status:   models.EnrollmentNotStarted,
	}

	if err := db.DB.Create(&enrollment).Error; err != nil {
		// Already enrolled - just return existing
		var existing models.Enrollment
		db.DB.Where("user_id = ? AND module_id = ?", userID, req.ModuleID).
			Preload("Module.Category").
			First(&existing)
		c.JSON(http.StatusOK, enrollmentToListJSON(&existing, false))
		return
	}

	db.DB.Preload("Module.Category").First(&enrollment, "id = ?", enrollment.ID)
	c.JSON(http.StatusCreated, enrollmentToListJSON(&enrollment, false))
}

func ptrString(s string) *string {
	return &s
}
