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
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id} [delete]
func (ac *Controller) DeleteLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	lobbyDeleted, err := ac.lobbyService.DeleteLobby(lobbyId, reqUser.Id)
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
	req := GetLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("GetLobby(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}

	lobby, err := ac.lobbyService.GetLobby(req.Id)
	if err != nil {
		// TODO: Move this into a generic wrapper at the router level
		fmt.Println(fmt.Errorf("GetLobby() error: %w", err))
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

// @Summary		Update Lobby
// @Description    Updates an existing lobby
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		lobby.UpdateLobbyRequest true	"Request Object"
// @Success		200
// @Router			/lobby/updateplayer [post]
func (ac *Controller) UpdateLobby(c *gin.Context) {
	req := UpdateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("UpdateLobby(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}
}

// @Summary		Add User to Lobby
// @Description    Add an existing User to an existing Lobby
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		lobby.AddUserToLobbyRequest true	"Request Object"
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id}/adduser [post]
func (ac *Controller) AddUserToLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	req := AddUserToLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("AddUserToLobby(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}

	// Load the lobby, so we can verify that the lobby exists
	lobby, err := ac.lobbyService.GetLobby(lobbyId)
	if err != nil {
		c.String(http.StatusInternalServerError, "Woops", err)
		return
	}
	if lobby == nil {
		c.String(http.StatusNotFound, "Lobby not found")
		return
	}

	// Verify that the request user is the owner of the lobby
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	if lobby.Owner != reqUser.Id {
		c.String(http.StatusUnauthorized, "Only the owner of the lobby can invite players")
		return
	}

	err = ac.lobbyService.AddUserToLobby(lobbyId, req.UserId)
	if err != nil {
		c.String(http.StatusBadRequest, "Unable to add user to lobby", err)
		return
	}
}
