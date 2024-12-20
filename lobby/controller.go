package lobby

import (
	"fmt"
	"main/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

type Controller struct {
	lobbyService *Service
}

func NewController(lobbyService *Service) *Controller {
	return &Controller{
		lobbyService: lobbyService,
	}
}

// @Summary		Creates a new lobby
// @Description	Creates a new lobby, with the requesting user as the owner of the lobby.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Produce		json
// @Param			request	body		lobby.CreateLobbyRequest	true	"Request Object"
// @Success		200		{object}	lobby.CreateLobbyResponse
// @Router			/lobby [post]
func (ac *Controller) CreateLobby(c *gin.Context) {
	req := CreateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	id, err := ac.lobbyService.CreateLobby(req.Name, reqUser.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create lobby. ", err)
	}
	resp := &CreateLobbyResponse{Id: id}

	c.JSON(http.StatusOK, resp)
}

// @Summary		Deletes a lobby
// @Description    Deletes a lobby.  The lobby must belong to the user sending the request.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		lobby.DeleteLobbyRequest	true	"Request Object"
// @Success		200
// @Router			/lobby [delete]
func (ac *Controller) DeleteLobby(c *gin.Context) {
	req := DeleteLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	lobbyDeleted, err := ac.lobbyService.DeleteLobby(req.Id, reqUser.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to delete lobby. ", err)
	}
	if !lobbyDeleted {
		c.String(http.StatusNotFound, "Unable to find lobby to delete")
	}
	c.Status(http.StatusOK)
}

// @Summary		Get a lobby
// @Description    Gets a lobby.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		lobby.GetLobbyRequest	true	"Request Object"
// @Success		200		{object}	lobby.GetLobbyResponse
// @Router			/lobby/getbyid [post]
func (ac *Controller) GetLobby(c *gin.Context) {
	// TODO: find a re-usable way to translate the context into a typed request
	req := GetLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("GetLobby(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}

	lobby, err := ac.lobbyService.GetLobby(req.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Woops", err)
		return
	}
	if lobby == nil {
		c.String(http.StatusNotFound, "Lobby not found")
		return
	}
	resp := GetLobbyResponse{Lobby: *lobby}
	c.JSON(http.StatusOK, resp)
}
