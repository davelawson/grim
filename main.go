package main

import (
	"database/sql"
	"fmt"
	"main/controller"
	"main/docs"
	"main/repo"
	"main/service"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	DATABASE_FILE := os.Getenv("GRIM_DB")
	fmt.Println("Database location: " + DATABASE_FILE)
	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		fmt.Println("Error opening database at {}: {}", DATABASE_FILE, err)
		return
	}
	GRIM_SSL := os.Getenv("GRIM_SSL")
	fmt.Println("SSL cert and key location: {}", GRIM_SSL)

	router := gin.Default()

	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	addUserRoutes(router, userController)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	fmt.Println("Listening for requests on port 8080")
	router.RunTLS("localhost:8080", GRIM_SSL+"/grim.crt", GRIM_SSL+"/grim.key")
}

func addSwaggerInfo() {
	docs.SwaggerInfo.Title = "Grimoire Backend API"
	docs.SwaggerInfo.Description = "This is the backend RESTful server for the Grimoire game."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/v2"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}

func addUserRoutes(router *gin.Engine, controller *controller.UserController) {
	group := router.Group("/user")
	group.POST("", controller.CreateUser)

	group = router.Group("/user/getbyemail")
	group.POST("", controller.GetUserByEmail)
}
