package controllers

import (
	"net/http"
	"pathfinder/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UserController holds the database dependency.
type UserController struct {
	DB *gorm.DB
}

// NewUserController creates a new UserController with the given database connection.
func NewUserController(db *gorm.DB) *UserController {
	return &UserController{DB: db}
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	var users []models.User

	// Find all users in the database
	result := uc.DB.Find(&users)

	// Handle potential errors during the query
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve users",
		})
		return
	}

	// Return the list of users
	ctx.JSON(http.StatusOK, gin.H{"users": users})
}
