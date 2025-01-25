package main

import (
	"pathfinder/controllers"
	"pathfinder/initializers"
	"pathfinder/middleware"

	"github.com/gin-gonic/gin"
)

func init() {
	initializers.LoadEnv()
	initializers.InitDB()
	initializers.SyncDb()
}

func main() {
	router := gin.Default()
	router.GET("/api/user", controllers.GetUsers)
	router.POST("/api/user", controllers.SignUp)
	router.POST("/api/login", controllers.Login)
	router.GET("/api/validate", middleware.RequireAuth, controllers.Validate)

	router.Run()

}
