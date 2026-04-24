package handlers

import (
	"net/http"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ListSlideNotes(c *gin.Context) {
	slideID := c.Param("id")

	var notes []models.Comment
	if err := db.DB.Where("entity_type = ? AND entity_id = ?", "slide", slideID).
		Preload("User").
		Order("created_at DESC").
		Find(&notes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notes"})
		return
	}
	out := make([]gin.H, 0, len(notes))
	for i := range notes {
		n := &notes[i]
		out = append(out, gin.H{
			"id":         n.ID,
			"userId":     n.UserID,
			"content":    n.Content,
			"entityType": n.EntityType,
			"entityId":   n.EntityID,
			"parentId":   n.ParentID,
			"createdAt":  n.CreatedAt,
			"updatedAt":  n.UpdatedAt,
			"user":       newAuthUserResponse(&n.User),
		})
	}
	c.JSON(http.StatusOK, out)
}

func CreateSlideNote(c *gin.Context) {
	slideID := c.Param("id")
	userID, ok := contextUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := models.Comment{
		ID:         uuid.New().String(),
		UserID:     userID,
		Content:    req.Content,
		EntityType: "slide",
		EntityID:   slideID,
	}

	if err := db.DB.Create(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
		return
	}

	db.DB.Preload("User").First(&note, "id = ?", note.ID)
	c.JSON(http.StatusCreated, gin.H{
		"id":         note.ID,
		"userId":     note.UserID,
		"content":    note.Content,
		"entityType": note.EntityType,
		"entityId":   note.EntityID,
		"parentId":   note.ParentID,
		"createdAt":  note.CreatedAt,
		"updatedAt":  note.UpdatedAt,
		"user":       newAuthUserResponse(&note.User),
	})
}

func UpdateNote(c *gin.Context) {
	noteID := c.Param("id")
	userID, uidOK := contextUserID(c)
	userRole, _ := contextUserRole(c)
	if !uidOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var note models.Comment
	if err := db.DB.First(&note, "id = ?", noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if note.UserID != userID && userRole != models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only edit your own notes"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	db.DB.Model(&note).Update("content", req.Content)
	db.DB.Preload("User").First(&note, "id = ?", noteID)
	c.JSON(http.StatusOK, gin.H{
		"id":         note.ID,
		"userId":     note.UserID,
		"content":    note.Content,
		"entityType": note.EntityType,
		"entityId":   note.EntityID,
		"parentId":   note.ParentID,
		"createdAt":  note.CreatedAt,
		"updatedAt":  note.UpdatedAt,
		"user":       newAuthUserResponse(&note.User),
	})
}

func DeleteNote(c *gin.Context) {
	noteID := c.Param("id")
	userID, uidOK := contextUserID(c)
	userRole, _ := contextUserRole(c)
	if !uidOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var note models.Comment
	if err := db.DB.First(&note, "id = ?", noteID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
		return
	}

	if note.UserID != userID && userRole != models.RoleSuperAdmin {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own notes"})
		return
	}

	if err := db.DB.Delete(&note).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Note deleted"})
}
