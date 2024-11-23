package controllers

import (
	"encoding/json"
	"main/models"
	"main/services"

	"github.com/gin-gonic/gin"
)

type albumService interface {
	GetAlbums() ([]models.Album, error)
}

type AlbumController struct {
	service *services.AlbumService
	bob     int
}

func NewAlbumController() *AlbumController {
	return &AlbumController{bob: 1}
}

func (ac *AlbumController) GetAlbums(c *gin.Context) {
	albums, err := ac.service.GetAlbums()
	if err != nil {
		// TODO: handle errors in some way here
		return
	}

	// TODO: this seems awful -- we need some way to standardize this!
	albumsJsonBytes, _ := json.Marshal(albums)
	albumsJsonString := string(albumsJsonBytes)
	c.JSON(200, gin.H{"albums": albumsJsonString})
}
