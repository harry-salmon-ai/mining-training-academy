package handlers

import (
	"net/http"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
)

func ListNotifications(c *gin.Context) {
	userID, _ := c.Get("userId")
	var notifications []models.Notification
	db.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(50).Find(&notifications)
	c.JSON(http.StatusOK, notifications)
}

func MarkNotificationRead(c *gin.Context) {
	id := c.Param("id")
	db.DB.Model(&models.Notification{}).Where("id = ?", id).Update("read", true)
	c.JSON(http.StatusOK, gin.H{"message": "Marked as read"})
}

func MarkAllNotificationsRead(c *gin.Context) {
	userID, _ := c.Get("userId")
	db.DB.Model(&models.Notification{}).Where("user_id = ? AND read = false", userID).Update("read", true)
	c.JSON(http.StatusOK, gin.H{"message": "All marked as read"})
}

func ListAuditLogs(c *gin.Context) {
	var logs []models.AuditLog
	query := db.DB.Preload("User").Order("created_at DESC").Limit(100)

	if action := c.Query("action"); action != "" {
		query = query.Where("action = ?", action)
	}
	if userID := c.Query("userId"); userID != "" {
		query = query.Where("user_id = ?", userID)
	}

	query.Find(&logs)
	c.JSON(http.StatusOK, logs)
}
