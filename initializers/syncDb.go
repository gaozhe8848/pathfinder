package initializers

import (
	"pathfinder/models"

	"gorm.io/gorm"
)

func SyncDb(db *gorm.DB) {
	db.AutoMigrate(&models.User{})
}
