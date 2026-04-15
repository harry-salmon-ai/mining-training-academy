package middleware

import (
	"net/http"

	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
)

type Permission struct {
	Resource string
	Action   string
}

var permissions = map[models.UserRole]map[string][]string{
	models.RoleSuperAdmin: {
		"modules":     {"create", "read", "update", "delete", "publish", "archive"},
		"users":       {"create", "read", "update", "delete", "assign_role"},
		"categories":  {"create", "read", "update", "delete"},
		"enrollments": {"create", "read", "update", "delete", "assign"},
		"reports":     {"read", "export"},
		"settings":    {"read", "update"},
		"audit":       {"read"},
	},
	models.RoleAdmin: {
		"modules":     {"create", "read", "update", "delete", "publish"},
		"users":       {"create", "read", "update"},
		"categories":  {"create", "read", "update"},
		"enrollments": {"create", "read", "update", "assign"},
		"reports":     {"read", "export"},
		"settings":    {"read"},
		"audit":       {"read"},
	},
	models.RoleInstructor: {
		"modules":     {"create", "read", "update"},
		"users":       {"read"},
		"categories":  {"read"},
		"enrollments": {"read", "assign"},
		"reports":     {"read"},
	},
	models.RoleSupervisor: {
		"modules":     {"read"},
		"users":       {"read"},
		"enrollments": {"read", "assign"},
		"reports":     {"read"},
	},
	models.RoleLearner: {
		"modules":     {"read"},
		"enrollments": {"read"},
	},
}

func CanAccess(role models.UserRole, resource, action string) bool {
	rolePerms, ok := permissions[role]
	if !ok {
		return false
	}
	actions, ok := rolePerms[resource]
	if !ok {
		return false
	}
	for _, a := range actions {
		if a == action {
			return true
		}
	}
	return false
}

func RequirePermission(resource, action string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userRole, ok := role.(models.UserRole)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid role"})
			c.Abort()
			return
		}

		if !CanAccess(userRole, resource, action) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}

func RequireAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("userRole")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}

		userRole, ok := role.(models.UserRole)
		if !ok || (userRole != models.RoleAdmin && userRole != models.RoleSuperAdmin) {
			c.JSON(http.StatusForbidden, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
