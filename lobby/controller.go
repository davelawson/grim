package lobby

import (
	"errors"
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
	lobbyService *ServiceFacade
}

func NewController(lobbyService *ServiceFacade) *Controller {
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
func (lc *Controller) CreateLobby(c *gin.Context) {
	req := CreateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body: %v", reqErr)
		return
	}
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	id, err := lc.lobbyService.CreateLobby(req.Name, reqUser.Id)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create lobby. %v", err)
	}
	resp := &CreateLobbyResponse{Id: *id}

	c.JSON(http.StatusOK, resp)
}

// @Summary		Deletes a lobby
// @Description    Deletes a lobby.  The lobby must belong to the user sending the request.
// @Security ApiKeyAuth
// @Tags			lobby
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id} [delete]
func (lc *Controller) DeleteLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	fmt.Println("DeleteLobby(): ", lobbyId)
	err := lc.lobbyService.DeleteLobby(lobbyId, reqUser.Id)
	if errors.Is(err, DeleteLobbyErrors.NotFound) {
		c.String(http.StatusNotFound, "Unable to process request: %v", err)
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Unable to process request: %v", err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Get a lobby
// @Description    Gets a lobby.
// @Security ApiKeyAuth
// @Tags			lobby
// @Param			id path string true "Lobby Id"
// @Success		200		{object}	lobby.GetLobbyResponse
// @Router			/lobby/{id} [get]
func (lc *Controller) GetLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	fmt.Println("GetLobby(): ", lobbyId)

	lobby, err := lc.lobbyService.GetLobby(lobbyId)
	if errors.Is(err, GetLobbyErrors.NotFound) {
		c.String(http.StatusNotFound, "Unable to process request: %v", err)
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Woops: %v", err)
		return
	}
	if lobby == nil {
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
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id} [put]
func (lc *Controller) UpdateLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	req := UpdateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("UpdateLobby() id: ", lobbyId, ", req: ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload: %v", reqErr)
		return
	}

	err := lc.lobbyService.UpdateLobby(lobbyId, req.Name, req.Owner)
	if err == UpdateLobbyErrors.NotFound {
		c.String(http.StatusNotFound, "Unable to process request: %v", err)
		return
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Woops: %v", err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Add User to Lobby
// @Description    Add an existing User to an existing Lobby
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		lobby.AddUserToLobbyRequest true	"Request Object"
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id}/user [post]
func (lc *Controller) AddUserToLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	req := AddUserToLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("AddUserToLobby() lobbyId: ", lobbyId, ", body: ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload: %v", reqErr)
		return
	}

	// Verify that the request user is the owner of the lobby
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)

	err := lc.lobbyService.AddUserToLobby(lobbyId, req.UserId, reqUser.Id)
	if err == AddUserToLobbyErrors.LobbyNotFound || err == AddUserToLobbyErrors.UserNotFound {
		c.String(http.StatusNotFound, "Unable to add user to lobby: %v", err)
	} else if err == AddUserToLobbyErrors.NotOwner {
		c.String(http.StatusForbidden, "Unable to add user to lobby: %v", err)
	} else if err == AddUserToLobbyErrors.UserAlreadyInLobby {
		c.String(http.StatusBadRequest, "Unable to add user to lobby: %v", err)
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Unable to add user to lobby: %v", err)
		return
	}
}

// @Summary		Remove User from Lobby
// @Description    Remove a lobby user from an existing Lobby
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			id path string true "Lobby Id"
// @Param			user_id path string true "User Id"
// @Success		200
// @Router			/lobby/{id}/user/{user_id} [delete]
func (lc *Controller) RemoveUserFromLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	userId := c.Param("user_id")
	fmt.Println("RemoveUserFromLobby() lobbyId: ", lobbyId, ", userId: ", userId)
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)

	err := lc.lobbyService.RemoveUserFromLobby(lobbyId, userId, reqUser.Id)
	if err == RemoveUserFromLobbyErrors.LobbyNotFound || err == RemoveUserFromLobbyErrors.UserNotInLobby {
		c.String(http.StatusNotFound, "Unable to remove user to lobby: %v", err)
	} else if err == RemoveUserFromLobbyErrors.NotOwner {
		c.String(http.StatusForbidden, "Unable to remove user to lobby: %v", err)
	} else if err != nil {
		c.String(http.StatusInternalServerError, "Unable to remove user to lobby: %v", err)
		return
	}
	c.Status(http.StatusOK)
}
