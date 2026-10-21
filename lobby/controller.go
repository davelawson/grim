package lobby

import (
	"errors"
	"fmt"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Controller struct {
	lobbyService *ServiceFacade
}

func NewController(lobbyService *ServiceFacade) *Controller {
	return &Controller{
		lobbyService: lobbyService,
	}
}

// @Summary		Creates a new lobby
// @Description Creates a public lobby with the requesting user as owner and member. Requires no existing open lobby membership.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Produce		json
// @Param			request	body		api.CreateLobbyRequest	true	"Request Object"
// @Success		200		{object}	api.CreateLobbyResponse
// @Failure 409,500 {object} api.ErrorResponse
// @Router			/lobby [post]
func (lc *Controller) CreateLobby(c *gin.Context) {
	req := api.CreateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body: %v", reqErr)
		return
	}
	reqUserAny, _ := c.Get("reqUser")
	reqUser := reqUserAny.(*model.User)
	id, err := lc.lobbyService.CreateLobby(req.Name, reqUser.Id)
	if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	resp := &api.CreateLobbyResponse{Id: *id}

	c.JSON(http.StatusOK, resp)
}

// @Summary		Deletes a lobby
// @Description    Soft-deletes an open lobby and retains its memberships. The lobby must belong to the user sending the request.
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
		util.WriteAPIError(c, err)
		return
	}

	c.Status(http.StatusOK)
}

// @Summary		Get a lobby
// @Description    Gets an open or closed lobby; authenticated members only.
// @Security ApiKeyAuth
// @Tags			lobby
// @Param			id path string true "Lobby Id"
// @Success		200		{object}	api.GetLobbyResponse
// @Router			/lobby/{id} [get]
func (lc *Controller) GetLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	fmt.Println("GetLobby(): ", lobbyId)

	reqUser := c.MustGet("reqUser").(*model.User)
	lobby, err := lc.lobbyService.GetLobby(lobbyId, reqUser.Id)
	if errors.Is(err, GetLobbyErrors.NotFound) {
		c.String(http.StatusNotFound, "Unable to process request: %v", err)
		return
	} else if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	resp := api.GetLobbyResponse{Lobby: *lobby}
	c.JSON(http.StatusOK, resp)
}

// @Summary		Update Lobby
// @Description    Current owner only; ownership transfers require a member target. Closed lobbies cannot be updated.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		api.UpdateLobbyRequest true	"Request Object"
// @Param			id path string true "Lobby Id"
// @Success		200
// @Router			/lobby/{id} [put]
func (lc *Controller) UpdateLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	req := api.UpdateLobbyRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("UpdateLobby() id: ", lobbyId, ", req: ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload: %v", reqErr)
		return
	}

	reqUser := c.MustGet("reqUser").(*model.User)
	err := lc.lobbyService.UpdateLobby(lobbyId, req.Name, req.Owner, reqUser.Id)
	if err == UpdateLobbyErrors.NotFound {
		c.String(http.StatusNotFound, "Unable to process request: %v", err)
		return
	} else if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// @Summary		Add User to Lobby
// @Description Owner only. Add an existing user to an open lobby with a free slot; the user must not belong to another open lobby.
// @Security ApiKeyAuth
// @Tags			lobby
// @Accept			json
// @Param			request	body		api.AddUserToLobbyRequest true	"Request Object"
// @Param			id path string true "Lobby Id"
// @Success		200
// @Failure 409,500 {object} api.ErrorResponse
// @Router			/lobby/{id}/user [post]
func (lc *Controller) AddUserToLobby(c *gin.Context) {
	lobbyId := c.Param("id")
	req := api.AddUserToLobbyRequest{}
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
		util.WriteAPIError(c, err)
		return
	}
	if err == nil {
		c.Status(http.StatusOK)
	}
}

// @Summary		Remove User from Lobby
// @Description    Owner-managed removal or voluntary departure. The owner must transfer ownership before leaving. Closed lobbies reject removal.
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
		return
	} else if err == RemoveUserFromLobbyErrors.NotOwner {
		c.String(http.StatusForbidden, "Unable to remove user to lobby: %v", err)
		return
	} else if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.Status(http.StatusOK)
}

// SetReady records only the authenticated member's own launch readiness.
// @Summary Set own lobby readiness
// @Security ApiKeyAuth
// @Tags lobby
// @Accept json
// @Produce json
// @Param id path string true "Lobby ID"
// @Param request body api.ReadyRequest true "Readiness (explicit true or false)"
// @Success 200 {object} api.GetLobbyResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/ready [put]
func (lc *Controller) SetReady(c *gin.Context) {
	var request api.ReadyRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteAPIError(c, util.ErrInvalidRequest)
		return
	}
	actor := c.MustGet("reqUser").(*model.User)
	lobby, err := lc.lobbyService.SetReady(c.Param("id"), actor.Id, *request.Ready)
	if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.JSON(http.StatusOK, api.GetLobbyResponse{Lobby: *lobby})
}

// Launch creates a persistent match in setup and permanently closes the lobby.
// @Summary Launch a lobby once
// @Description Owner only, with two to four total participants and all humans ready. Automated opponents draft immediately. Identical request-ID retries replay the original response while the match is active, or return 404 after deletion; a new ID on a closed lobby conflicts.
// @Security ApiKeyAuth
// @Tags lobby
// @Accept json
// @Produce json
// @Param id path string true "Lobby ID"
// @Param request body api.LaunchRequest true "Launch request ID"
// @Success 201 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/launch [post]
func (lc *Controller) Launch(c *gin.Context) {
	var request api.LaunchRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		util.WriteAPIError(c, util.ErrInvalidRequest)
		return
	}
	actor := c.MustGet("reqUser").(*model.User)
	view, err := lc.lobbyService.Launch(c.Param("id"), actor.Id, request.RequestID)
	if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.JSON(http.StatusCreated, api.MatchResponse{Match: view})
}

// Join adds the authenticated user to a public open lobby.
// @Summary Join a public lobby
// @Description Requires a free slot and no other open lobby membership. Repeated joins return the current lobby without clearing readiness.
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Param id path string true "Lobby ID"
// @Success 200 {object} api.GetLobbyResponse
// @Failure 401,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/join [post]
func (lc *Controller) Join(c *gin.Context) {
	actor := c.MustGet("reqUser").(*model.User)
	lobby, err := lc.lobbyService.Join(c.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.JSON(http.StatusOK, api.GetLobbyResponse{Lobby: *lobby})
}

// Leave removes only the authenticated user using the transactional removal flow.
// @Summary Leave a lobby
// @Description The owner must transfer ownership before leaving. Departure clears remaining members' readiness. Closed or deleted lobbies reject departure.
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Param id path string true "Lobby ID"
// @Success 200
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/leave [delete]
func (lc *Controller) Leave(c *gin.Context) {
	actor := c.MustGet("reqUser").(*model.User)
	err := lc.lobbyService.RemoveUserFromLobby(c.Param("id"), actor.Id, actor.Id)
	switch {
	case errors.Is(err, RemoveUserFromLobbyErrors.LobbyNotFound), errors.Is(err, RemoveUserFromLobbyErrors.UserNotInLobby):
		err = util.ErrNotFound
	case errors.Is(err, RemoveUserFromLobbyErrors.NotOwner):
		err = util.ErrForbidden
	}
	if err != nil {
		util.WriteAPIError(c, err)
		return
	}
	c.Status(http.StatusOK)
}
