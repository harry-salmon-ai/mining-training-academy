package handlers

import (
	"math"
	"net/http"
	"time"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetDashboardStats(c *gin.Context) {
	userRole, ok := contextUserRole(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if userRole == models.RoleAdmin || userRole == models.RoleSuperAdmin {
		getAdminDashboard(c)
	} else {
		getLearnerDashboard(c)
	}
}

func getAdminDashboard(c *gin.Context) {
	var totalUsers, activeUsers, totalModules, publishedModules int64
	var totalEnrollments, completedEnrollments int64

	thirtyDaysAgo := time.Now().AddDate(0, 0, -30)

	db.DB.Model(&models.User{}).Where("is_active = true").Count(&totalUsers)
	db.DB.Model(&models.User{}).Where("last_login_at > ?", thirtyDaysAgo).Count(&activeUsers)
	db.DB.Model(&models.Module{}).Count(&totalModules)
	db.DB.Model(&models.Module{}).Where("status = ?", models.StatusPublished).Count(&publishedModules)
	db.DB.Model(&models.Enrollment{}).Count(&totalEnrollments)
	db.DB.Model(&models.Enrollment{}).Where("status = ?", models.EnrollmentCompleted).Count(&completedEnrollments)

	completionRate := 0.0
	if totalEnrollments > 0 {
		completionRate = math.Round(float64(completedEnrollments) / float64(totalEnrollments) * 100)
	}

	type AvgResult struct {
		Avg *float64
	}
	var avgResult AvgResult
	db.DB.Model(&models.QuizAttempt{}).Select("AVG(score) as avg").Where("score IS NOT NULL").Scan(&avgResult)
	avgScore := 0.0
	if avgResult.Avg != nil {
		avgScore = math.Round(*avgResult.Avg)
	}

	var recentActivity []models.AuditLog
	db.DB.Preload("User").Order("created_at DESC").Limit(20).Find(&recentActivity)

	type auditLogOut struct {
		ID        string           `json:"id"`
		UserID    string           `json:"userId"`
		Action    string           `json:"action"`
		Entity    string           `json:"entity"`
		EntityID  *string          `json:"entityId"`
		Metadata  datatypes.JSON   `json:"metadata"`
		IPAddress *string          `json:"ipAddress"`
		CreatedAt time.Time        `json:"createdAt"`
		User      authUserResponse `json:"user,omitempty"`
	}
	recentOut := make([]auditLogOut, 0, len(recentActivity))
	for _, log := range recentActivity {
		recentOut = append(recentOut, auditLogOut{
			ID:        log.ID,
			UserID:    log.UserID,
			Action:    log.Action,
			Entity:    log.Entity,
			EntityID:  log.EntityID,
			Metadata:  log.Metadata,
			IPAddress: log.IPAddress,
			CreatedAt: log.CreatedAt,
			User:      newAuthUserResponse(&log.User),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"totalUsers":           totalUsers,
		"activeUsers":          activeUsers,
		"totalModules":         totalModules,
		"publishedModules":     publishedModules,
		"totalEnrollments":     totalEnrollments,
		"completedEnrollments": completedEnrollments,
		"completionRate":       completionRate,
		"avgScore":             avgScore,
		"recentActivity":       recentOut,
	})
}

func getLearnerDashboard(c *gin.Context) {
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var enrollmentCount, completedCount int64
	db.DB.Model(&models.Enrollment{}).Where("user_id = ?", userID).Count(&enrollmentCount)
	db.DB.Model(&models.Enrollment{}).Where("user_id = ? AND status = ?", userID, models.EnrollmentCompleted).Count(&completedCount)

	var inProgress []models.Enrollment
	db.DB.Where("user_id = ? AND status = ?", userID, models.EnrollmentInProgress).
		Preload("Module.Category").
		Order("updated_at DESC").
		Find(&inProgress)

	var certCount int64
	db.DB.Model(&models.Certificate{}).Where("user_id = ?", userID).Count(&certCount)

	type AvgResult struct {
		Avg *float64
	}
	var avgResult AvgResult
	db.DB.Model(&models.QuizAttempt{}).Select("AVG(score) as avg").
		Where("user_id = ? AND passed = true", userID).Scan(&avgResult)
	avgScore := 0.0
	if avgResult.Avg != nil {
		avgScore = math.Round(*avgResult.Avg)
	}

	inProgressOut := make([]gin.H, 0, len(inProgress))
	for _, e := range inProgress {
		row := gin.H{
			"id":       e.ID,
			"status":   e.Status,
			"progress": e.Progress,
		}
		mod := gin.H{"title": e.Module.Title}
		if e.Module.CategoryID != "" {
			mod["category"] = gin.H{"name": e.Module.Category.Name}
		}
		row["module"] = mod
		inProgressOut = append(inProgressOut, row)
	}

	c.JSON(http.StatusOK, gin.H{
		"enrollments":      enrollmentCount,
		"completedModules": completedCount,
		"inProgress":       inProgressOut,
		"certificates":     certCount,
		"avgScore":         avgScore,
	})
}
