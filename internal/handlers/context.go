package handlers

import (
	"github.com/gin-gonic/gin"
	"go-backend-react-frontend/internal/models"
)

func contextUserRole(c *gin.Context) (models.UserRole, bool) {
	v, ok := c.Get("userRole")
	if !ok || v == nil {
		return "", false
	}
	switch r := v.(type) {
	case models.UserRole:
		return r, true
	case string:
		if r == "" {
			return "", false
		}
		return models.UserRole(r), true
	default:
		return "", false
	}
}

func contextUserID(c *gin.Context) (string, bool) {
	v, ok := c.Get("userId")
	if !ok || v == nil {
		return "", false
	}
	s, ok := v.(string)
	return s, ok && s != ""
}
