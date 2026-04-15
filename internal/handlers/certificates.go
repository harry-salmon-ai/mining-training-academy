package handlers

import (
	"fmt"
	"net/http"
	"time"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListCertificates(c *gin.Context) {
	userID, _ := c.Get("userId")
	role, _ := c.Get("userRole")
	userRole := role.(models.UserRole)

	query := db.DB.Preload("Module").Preload("User")

	if userRole == models.RoleLearner {
		query = query.Where("user_id = ?", userID)
	} else if uid := c.Query("userId"); uid != "" {
		query = query.Where("user_id = ?", uid)
	}

	var certs []models.Certificate
	query.Order("issued_at DESC").Find(&certs)
	c.JSON(http.StatusOK, certs)
}

func GenerateCertificate(c *gin.Context) {
	var req struct {
		UserID   string   `json:"userId" binding:"required"`
		ModuleID string   `json:"moduleId" binding:"required"`
		Score    *float64 `json:"score"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	certNo := fmt.Sprintf("MTA-%s-%d", req.ModuleID[:8], time.Now().UnixMilli())

	cert := models.Certificate{
		ID:            uuid.New().String(),
		UserID:        req.UserID,
		ModuleID:      req.ModuleID,
		CertificateNo: certNo,
		Score:         req.Score,
	}

	if err := db.DB.Create(&cert).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate certificate"})
		return
	}

	db.DB.Preload("Module").Preload("User").First(&cert, "id = ?", cert.ID)
	c.JSON(http.StatusCreated, cert)
}
