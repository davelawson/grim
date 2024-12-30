package auth

import (
	"fmt"
	"main/model/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Controller struct {
	authService *ServiceFacade
}

func NewController(authService *ServiceFacade) *Controller {
	return &Controller{authService: authService}
}

// Login godoc
//
//	@Summary		Authenticate user and return bearer token
//	@Description	Verifies the user password, and generates a new bearer token for that user.
//	@Tags			login
//	@Accept			json
//	@Produce		json
//	@Param			request	body		api.LoginRequest	true	"Request Object"
//	@Success		200		{object}	api.LoginResponse
//	@Router			/login [post]
func (ac *Controller) Login(c *gin.Context) {
	req := api.LoginRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("Login(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	req.Email = strings.ToLower(req.Email)
	bearerToken, authErr := ac.authService.Login(req.Email, req.Password)
	if authErr != nil {
		fmt.Println("InternalServerError detected: ", authErr)
		c.String(http.StatusInternalServerError, "Something went wrong!")
		return
	}
	if bearerToken == nil {
		c.String(http.StatusUnauthorized, "Bad user name or password")
		return
	}
	resp := &api.LoginResponse{Token: *bearerToken}
	c.JSON(http.StatusOK, resp)
}
