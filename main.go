package main

import (
	"log"
	"os"
	"pathfinder/controllers"
	"pathfinder/initializers"
	"pathfinder/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize application
	db, err := initializers.InitDB()
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}
	jwtSecret := os.Getenv("SECRET")
	if jwtSecret == "" {
		log.Fatal("SECRET environment variable not set")
	}
	// ginMode := os.Getenv("GIN_MODE") // No longer needed for AuthController
	initializers.SyncDb(db)

	// Initialize controllers with the DB instance
	userController := controllers.NewUserController(db)
	authController := controllers.NewAuthController(db, jwtSecret)

	// Set up router
	router := gin.Default()

	// Setup routes
	// Protect the GetUsers route with RequireAuth middleware
	router.GET("/api/user", middleware.RequireAuth(db, jwtSecret), userController.GetUsers)
	router.POST("/api/user", authController.SignUp)
	router.POST("/api/login", authController.Login)
	// Pass the db instance and jwtSecret to the RequireAuth middleware factory
	router.GET("/api/validate", middleware.RequireAuth(db, jwtSecret), authController.Validate)

	router.Run()
}
