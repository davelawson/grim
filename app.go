package main

import (
	"database/sql"
	"fmt"
	"main/controllers"
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

	results, err := db.Query("select * from players")
	if err != nil {
		fmt.Println("Error executing query: {}", DATABASE_FILE, err)
		return
	}
	fmt.Println("Results from database: {}", results)

	router := gin.Default()

	albumController := controllers.NewAlbumController()
	addAlbumRoutes(router, albumController)

	router.Run("localhost:8080")
}

func addAlbumRoutes(router *gin.Engine, controller *controllers.AlbumController) {
	router.GET("/albums", func(c *gin.Context) {
		controller.GetAlbums(c)
	})
}
