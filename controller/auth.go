package controller

import (
	"fmt"
	"main/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type AuthController struct {
	authService *service.AuthService
}

type loginRequest struct {
	Email    string
	Password string
}

type loginResponse struct {
	Token []byte
}

func NewAuthController(authService *service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

// Login godoc
//
//	@Summary		Authenticate user and return bearer token
//	@Description	Verifies the user password, and generates a new bearer token for that user.
//	@Tags			login
//	@Accept			json
//	@Produce		json
//	@Param			request	body		loginRequest	true	"Request Object"
//	@Success		200		{object}	loginResponse
//	@Router			/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	req := loginRequest{}
	fmt.Println("Login(): {}", req)
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid JSON")
		return
	}

	bearerToken, authErr := ac.authService.Login(req.Email, req.Password)
	if authErr != nil {
		fmt.Println("InternalServerError detected: {}", authErr)
		c.String(http.StatusInternalServerError, "Something went wrong!")
		return
	}
	if bearerToken == nil {
		c.String(http.StatusUnauthorized, "Bad user name or password")
		return
	}
	resp := &loginResponse{Token: bearerToken}
	c.JSON(http.StatusOK, resp)
}
