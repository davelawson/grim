package main

import (
	"database/sql"
	"fmt"
	"main/controller"
	"main/repo"
	"main/service"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	DATABASE_FILE := os.Getenv("GRIM_DB")
	fmt.Println("Database location: " + DATABASE_FILE)
	db, err := sql.Open("sqlite3", DATABASE_FILE)
	if err != nil {
		fmt.Println("Error opening database at {}: {}", DATABASE_FILE, err)
		return
	}

	router := gin.Default()

	albumController := controller.NewAlbumController()
	addAlbumRoutes(router, albumController)

	userRepo := repo.NewUserRepo(db)
	userService := service.NewUserService(userRepo)
	userController := controller.NewUserController(userService)
	addUserRoutes(router, userController)

	router.Run("localhost:8080")
}

func addAlbumRoutes(router *gin.Engine, controller *controller.AlbumController) {
	router.GET("/albums", func(c *gin.Context) {
		controller.GetAlbums(c)
	})
}

func addUserRoutes(router *gin.Engine, controller *controller.UserController) {
	router.GET("/user", func(c *gin.Context) {
		controller.GetUserByEmail(c)
	})
	router.POST("/user", func(c *gin.Context) {
		controller.CreateUser(c)
	})
}
