package main

import (
	middleware "golang-training/task2/auth"
	"golang-training/task2/controller"
	"golang-training/task2/dao"
	"golang-training/task2/service"

	"golang-training/task2/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           User Management API
// @version         1.0
// @description     REST API for user authentication and management
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	r := gin.Default()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	docs.SwaggerInfo.Title = "User Management API"
	docs.SwaggerInfo.Description = "REST API for user authentication and management with JWT"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	daoImpl := &dao.MongoDAO{}
	userService := service.NewUserService(daoImpl)
	userController := controller.NewController(userService)

	auth := r.Group("")
	{
		auth.POST("/login", userController.Login)
		auth.POST("/logout", middleware.AuthMiddleware(userService), userController.Logout)
	}

	users := r.Group("/users")
	{
		users.POST("", userController.RegisterUser)
		users.GET("/me", middleware.AuthMiddleware(userService), userController.GetCurrentUser)
		users.GET("/:id", middleware.AuthMiddleware(userService), userController.GetUserByID)
		users.GET("", middleware.AuthMiddleware(userService), userController.GetAllUsers)
	}

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
