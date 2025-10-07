package main

import (
	middleware "golang-training/task2/auth"
	"golang-training/task2/controller"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	auth := r.Group("")
	{
		auth.POST("/login", controller.Login)
		auth.POST("/logout", middleware.AuthMiddleware(), controller.Logout)
	}

	users := r.Group("/users")
	{
		users.POST("", controller.RegisterUser)
		users.GET("/me", middleware.AuthMiddleware(), controller.GetCurrentUser)
		users.GET("/:id", middleware.AuthMiddleware(), controller.GetUserByID)
		users.GET("", middleware.AuthMiddleware(), controller.GetAllUsers)
	}

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
