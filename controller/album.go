package controller

import (
	"encoding/json"
	"main/model"
	"main/service"

	"github.com/gin-gonic/gin"
)

type albumService interface {
	GetAlbums() ([]model.Album, error)
}

type AlbumController struct {
	service *service.AlbumService
}

func NewAlbumController() *AlbumController {
	return &AlbumController{}
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
