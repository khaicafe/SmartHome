package controllers

import (
	"go-react-app/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func GetSettings(c *gin.Context) {
	var settings []models.Setting
	models.DB.Find(&settings)
	c.JSON(http.StatusOK, settings)
}

// POST /api/settings
func UpdateSetting(c *gin.Context) {
	var updates []models.Setting
	if err := c.ShouldBindJSON(&updates); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	for _, s := range updates {
		models.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			UpdateAll: true,
		}).Create(&s)
	}

	c.JSON(http.StatusOK, gin.H{"message": "Settings updated"})
}
