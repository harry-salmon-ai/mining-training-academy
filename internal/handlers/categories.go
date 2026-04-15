package handlers

import (
	"net/http"
	"strings"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListCategories(c *gin.Context) {
	var categories []models.Category
	query := db.DB.Order("sort_order ASC")

	if c.Query("active") == "true" {
		query = query.Where("is_active = true")
	}

	if err := query.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	c.JSON(http.StatusOK, categories)
}

func CreateCategory(c *gin.Context) {
	var req struct {
		Name        string  `json:"name" binding:"required"`
		Description *string `json:"description"`
		Color       string  `json:"color"`
		Icon        string  `json:"icon"`
		SortOrder   int     `json:"sortOrder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	cat := models.Category{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		Color:       req.Color,
		Icon:        req.Icon,
		SortOrder:   req.SortOrder,
		IsActive:    true,
	}
	if cat.Color == "" {
		cat.Color = "#185FA5"
	}
	if cat.Icon == "" {
		cat.Icon = "M"
	}

	if err := db.DB.Create(&cat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create category"})
		return
	}
	c.JSON(http.StatusCreated, cat)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	if err := db.DB.First(&cat, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
		Color       *string `json:"color"`
		Icon        *string `json:"icon"`
		SortOrder   *int    `json:"sortOrder"`
		IsActive    *bool   `json:"isActive"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updates := map[string]interface{}{}
	if req.Name != nil {
		updates["name"] = *req.Name
		updates["slug"] = strings.ToLower(strings.ReplaceAll(*req.Name, " ", "-"))
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Icon != nil {
		updates["icon"] = *req.Icon
	}
	if req.SortOrder != nil {
		updates["sort_order"] = *req.SortOrder
	}
	if req.IsActive != nil {
		updates["is_active"] = *req.IsActive
	}

	db.DB.Model(&cat).Updates(updates)
	db.DB.First(&cat, "id = ?", id)
	c.JSON(http.StatusOK, cat)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Category{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete category"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
