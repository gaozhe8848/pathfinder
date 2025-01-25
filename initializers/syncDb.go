package initializers

import "pathfinder/models"

func SyncDb() {
	DB.AutoMigrate(&models.User{})
}
