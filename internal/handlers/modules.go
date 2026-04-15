package handlers

import (
	"encoding/json"
	"net/http"

	"gorm.io/datatypes"
	"strconv"
	"strings"
	"time"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func ListModules(c *gin.Context) {
	role, _ := c.Get("userRole")
	userRole, _ := role.(models.UserRole)

	query := db.DB.Preload("Category").Preload("Author").Preload("Tags")

	if userRole == models.RoleLearner || userRole == "" {
		query = query.Where("status = ?", models.StatusPublished)
	} else if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", strings.ToUpper(status))
	}

	if catID := c.Query("category"); catID != "" {
		query = query.Where("category_id = ?", catID)
	}
	if level := c.Query("level"); level != "" {
		query = query.Where("level = ?", strings.ToUpper(level))
	}
	if search := c.Query("search"); search != "" {
		like := "%" + search + "%"
		query = query.Where("title ILIKE ? OR description ILIKE ?", like, like)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	var total int64
	query.Model(&models.Module{}).Count(&total)

	var modules []models.Module
	query.Order("updated_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&modules)

	c.JSON(http.StatusOK, gin.H{
		"modules": modules,
		"pagination": gin.H{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

func GetModule(c *gin.Context) {
	id := c.Param("id")
	var module models.Module

	if err := db.DB.Preload("Category").Preload("Author").Preload("Tags").
		Preload("Sections", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort_order ASC")
		}).
		Preload("Sections.Slides", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort_order ASC")
		}).
		Where("slug = ? OR id = ?", id, id).First(&module).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	c.JSON(http.StatusOK, module)
}

func CreateModule(c *gin.Context) {
	var req struct {
		Title       string             `json:"title" binding:"required"`
		Description *string            `json:"description"`
		CategoryID  string             `json:"categoryId" binding:"required"`
		Level       models.ModuleLevel `json:"level"`
		Duration    *int               `json:"duration"`
		Tags        []string           `json:"tags"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authorID, _ := c.Get("userId")
	slug := strings.ToLower(strings.ReplaceAll(strings.TrimSpace(req.Title), " ", "-"))
	slug = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return -1
	}, slug)

	if req.Level == "" {
		req.Level = models.LevelFoundation
	}

	module := models.Module{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Slug:        slug,
		Description: req.Description,
		Status:      models.StatusDraft,
		Level:       req.Level,
		Duration:    req.Duration,
		CategoryID:  req.CategoryID,
		AuthorID:    authorID.(string),
	}

	if err := db.DB.Create(&module).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create module"})
		return
	}

	for _, tag := range req.Tags {
		db.DB.Create(&models.ModuleTag{
			ID:       uuid.New().String(),
			ModuleID: module.ID,
			Tag:      tag,
		})
	}

	db.DB.Preload("Category").Preload("Tags").First(&module, "id = ?", module.ID)
	c.JSON(http.StatusCreated, module)
}

func UpdateModule(c *gin.Context) {
	id := c.Param("id")
	var module models.Module
	if err := db.DB.First(&module, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Module not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&module).Updates(req)
	db.DB.Preload("Category").Preload("Author").Preload("Tags").First(&module, "id = ?", id)
	c.JSON(http.StatusOK, module)
}

func DeleteModule(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Delete(&models.Module{}, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete module"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Module deleted"})
}

func PublishModule(c *gin.Context) {
	id := c.Param("id")
	now := time.Now()
	if err := db.DB.Model(&models.Module{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":       models.StatusPublished,
		"published_at": now,
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to publish module"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Module published"})
}

func ArchiveModule(c *gin.Context) {
	id := c.Param("id")
	if err := db.DB.Model(&models.Module{}).Where("id = ?", id).Update("status", models.StatusArchived).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to archive module"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Module archived"})
}

// Sections

func ListSections(c *gin.Context) {
	moduleID := c.Param("id")
	var sections []models.Section
	db.DB.Where("module_id = ?", moduleID).
		Preload("Slides", func(tx *gorm.DB) *gorm.DB {
			return tx.Order("sort_order ASC")
		}).
		Order("sort_order ASC").Find(&sections)
	c.JSON(http.StatusOK, sections)
}

func CreateSection(c *gin.Context) {
	moduleID := c.Param("id")
	var req struct {
		Title       string  `json:"title" binding:"required"`
		Description *string `json:"description"`
		SortOrder   int     `json:"sortOrder"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	section := models.Section{
		ID:          uuid.New().String(),
		Title:       req.Title,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		ModuleID:    moduleID,
	}

	if err := db.DB.Create(&section).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create section"})
		return
	}
	c.JSON(http.StatusCreated, section)
}

func UpdateSection(c *gin.Context) {
	id := c.Param("id")
	var section models.Section
	if err := db.DB.First(&section, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Section not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Model(&section).Updates(req)
	c.JSON(http.StatusOK, section)
}

func DeleteSection(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Section{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "Section deleted"})
}

// Slides

func CreateSlide(c *gin.Context) {
	sectionID := c.Param("id")
	var req struct {
		Title     string           `json:"title" binding:"required"`
		Type      models.SlideType `json:"type"`
		Content   json.RawMessage  `json:"content"`
		SortOrder int              `json:"sortOrder"`
		Notes     *string          `json:"notes"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Type == "" {
		req.Type = models.SlideContent
	}

	slide := models.Slide{
		ID:        uuid.New().String(),
		Title:     req.Title,
		Type:      req.Type,
		Content:   datatypes.JSON(req.Content),
		SortOrder: req.SortOrder,
		Notes:     req.Notes,
		SectionID: sectionID,
	}

	if err := db.DB.Create(&slide).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create slide"})
		return
	}
	c.JSON(http.StatusCreated, slide)
}

func GetSlide(c *gin.Context) {
	id := c.Param("id")
	var slide models.Slide
	if err := db.DB.Preload("Media").First(&slide, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found"})
		return
	}
	c.JSON(http.StatusOK, slide)
}

func UpdateSlide(c *gin.Context) {
	id := c.Param("id")
	var slide models.Slide
	if err := db.DB.First(&slide, "id = ?", id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Slide not found"})
		return
	}

	var req map[string]interface{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	db.DB.Model(&slide).Updates(req)
	db.DB.First(&slide, "id = ?", id)
	c.JSON(http.StatusOK, slide)
}

func DeleteSlide(c *gin.Context) {
	id := c.Param("id")
	db.DB.Delete(&models.Slide{}, "id = ?", id)
	c.JSON(http.StatusOK, gin.H{"message": "Slide deleted"})
}
