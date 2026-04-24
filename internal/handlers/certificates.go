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

func certificateToListJSON(c *models.Certificate, includeUser bool) gin.H {
	h := gin.H{
		"id":            c.ID,
		"userId":        c.UserID,
		"moduleId":      c.ModuleID,
		"certificateNo": c.CertificateNo,
		"issuedAt":      c.IssuedAt,
		"expiresAt":     c.ExpiresAt,
		"score":         c.Score,
		"pdfUrl":        c.PDFURL,
		"module": gin.H{
			"id":    c.Module.ID,
			"title": c.Module.Title,
			"slug":  c.Module.Slug,
		},
	}
	if includeUser {
		h["user"] = newAuthUserResponse(&c.User)
	}
	return h
}

func ListCertificates(c *gin.Context) {
	userID, uidOK := contextUserID(c)
	userRole, roleOK := contextUserRole(c)
	if !uidOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	if !roleOK {
		userRole = models.RoleLearner
	}

	query := db.DB.Preload("Module").Preload("User")

	if userRole == models.RoleLearner {
		query = query.Where("user_id = ?", userID)
	} else if uid := c.Query("userId"); uid != "" {
		query = query.Where("user_id = ?", uid)
	}

	var certs []models.Certificate
	query.Order("issued_at DESC").Find(&certs)
	includeUser := userRole != models.RoleLearner
	out := make([]gin.H, 0, len(certs))
	for i := range certs {
		out = append(out, certificateToListJSON(&certs[i], includeUser))
	}
	c.JSON(http.StatusOK, out)
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
