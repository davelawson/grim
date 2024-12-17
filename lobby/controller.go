package lobby

import (
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

// CreateLobby godoc
//
//	@Summary		Creates a new lobby
//	@Description	Creates a new lobby, with the requesting user as the owner of the lobby.
//	@Security ApiKeyAuth
//	@Tags			lobby
//	@Accept			json
//	@Produce		json
//	@Param			request	body		lobby.CreateLobbyRequest	true	"Request Object"
//	@Success		200		{object}	lobby.CreateLobbyResponse
//	@Router			/lobby [post]
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

func (ac *Controller) DeleteLobby(c *gin.Context) {
}

func (ac *Controller) GetLobby(c *gin.Context) {
}
