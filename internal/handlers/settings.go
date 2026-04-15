package handlers

import (
	"net/http"

	"go-backend-react-frontend/internal/db"
	"go-backend-react-frontend/internal/models"

	"github.com/gin-gonic/gin"
)

func GetSettings(c *gin.Context) {
	var settings models.PlatformSettings
	if err := db.DB.First(&settings, "id = ?", "default").Error; err != nil {
		settings = models.PlatformSettings{
			ID:           "default",
			PlatformName: "Mining Training Academy",
			Tagline:      "Operational Excellence Through Learning",
			PrimaryColor: "#1a1a1a",
			Timezone:     "UTC",
		}
		db.DB.Create(&settings)
	}
	c.JSON(http.StatusOK, settings)
}

func UpdateSettings(c *gin.Context) {
	var settings models.PlatformSettings
	db.DB.First(&settings, "id = ?", "default")

	var req map[string]interface{}
	c.ShouldBindJSON(&req)
	delete(req, "id")

	db.DB.Model(&settings).Updates(req)
	db.DB.First(&settings, "id = ?", "default")
	c.JSON(http.StatusOK, settings)
}
