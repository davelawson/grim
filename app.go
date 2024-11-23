package main

import (
	"main/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
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
