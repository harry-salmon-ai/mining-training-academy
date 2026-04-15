package handlers

import (
	"net/http"
	"strconv"

	"go-backend-react-frontend/internal/auth"
	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}

	query := db.DB.Model(&models.User{})

	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", like, like)
	}
	if role := c.Query("role"); role != "" {
		query = query.Where("role = ?", role)
	}
	if site := c.Query("siteLocation"); site != "" {
		query = query.Where("site_location = ?", site)
	}

	var total int64
	query.Count(&total)

	var users []models.User
	query.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"users": users,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func CreateUser(c *gin.Context) {
	var req struct {
		Email        string          `json:"email" binding:"required,email"`
		FirstName    string          `json:"firstName" binding:"required"`
		LastName     string          `json:"lastName" binding:"required"`
		Role         models.UserRole `json:"role"`
		Department   *string         `json:"department"`
		JobTitle     *string         `json:"jobTitle"`
		SiteLocation *string         `json:"siteLocation"`
		EmployeeID   *string         `json:"employeeId"`
		Password     string          `json:"password" binding:"required,min=8"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Role == "" {
		req.Role = models.RoleLearner
	}

	hash, _ := auth.HashPassword(req.Password)
	name := req.FirstName + " " + req.LastName

	user := models.User{
		ID:           uuid.New().String(),
		Email:        req.Email,
		Name:         &name,
		FirstName:    &req.FirstName,
		LastName:     &req.LastName,
		PasswordHash: &hash,
		Role:         req.Role,
		Department:   req.Department,
		JobTitle:     req.JobTitle,
		SiteLocation: req.SiteLocation,
		EmployeeID:   req.EmployeeID,
		IsActive:     true,
	}

	if err := db.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := db.DB.First(&user, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "passwordHash")
	delete(req, "id")

	db.DB.Model(&user).Updates(req)
	db.DB.First(&user, "id = ?", id)
	c.JSON(http.StatusOK, user)
}

func DeactivateUser(c *gin.Context) {
	id := c.Param("id")
	db.DB.Model(&models.User{}).Where("id = ?", id).Update("is_active", false)
	c.JSON(http.StatusOK, gin.H{"message": "User deactivated"})
}

func UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userId")
	var user models.User
	if err := db.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	var req struct {
		FirstName    *string `json:"firstName"`
		LastName     *string `json:"lastName"`
		Department   *string `json:"department"`
		JobTitle     *string `json:"jobTitle"`
		SiteLocation *string `json:"siteLocation"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.FirstName != nil && req.LastName != nil {
		name := *req.FirstName + " " + *req.LastName
		updates["name"] = name
	}
	if req.Department != nil {
		updates["department"] = *req.Department
	}
	if req.JobTitle != nil {
		updates["job_title"] = *req.JobTitle
	}
	if req.SiteLocation != nil {
		updates["site_location"] = *req.SiteLocation
	}

	db.DB.Model(&user).Updates(updates)
	db.DB.First(&user, "id = ?", userID)
	c.JSON(http.StatusOK, user)
}
