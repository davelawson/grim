package user

import (
	"fmt"
	"main/model"
	"main/model/api"
	"main/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

type Controller struct {
	userService *Service
	authService authService
}

func NewController(userService *Service, authService authService) *Controller {
	return &Controller{
		userService: userService,
		authService: authService,
	}
}

// GetUserByEmail godoc
//
//	@Summary		Get user by email
//	@Description	Lookup a specific user by email
//	@Security ApiKeyAuth
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body		api.GetUserRequest	true	"Request Object"
//	@Success		200		{object}	api.GetUserResponse
//	@Router			/user/getbyemail [post]
func (us *Controller) GetUserByEmail(c *gin.Context) {
	// TODO: Extract this auth block to a reusable chunk of code
	reqUser, authErr := us.authService.VerifyBearerToken(util.GetBearerToken(c))
	if authErr != nil {
		c.String(http.StatusInternalServerError, "Unable to verify validity of authentication token", authErr)
		return
	}
	if reqUser == nil {
		c.String(http.StatusUnauthorized, "Bad or missing authentication token")
		return
	}
	// TODO: find a re-usable way to translate the context into a typed request
	req := api.GetUserRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("GetUserByEmail(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}

	user, err := us.userService.GetUserByEmail(req.Email)
	if err != nil {
		c.String(http.StatusInternalServerError, "Woops", err)
		return
	}
	if user == nil {
		c.String(http.StatusNotFound, "User not found")
		return
	}

	resp := api.GetUserResponse{User: api.NewUser(user)}

	c.JSON(http.StatusOK, resp)
}

// CreateUser godoc
//
//	@Summary		Create user
//	@Description	Create a new user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			request	body	api.CreateUserRequest	true	"Request Object"
//	@Success		200
//	@Router			/user [post]
func (us *Controller) CreateUser(c *gin.Context) {
	req := api.CreateUserRequest{}
	fmt.Println("CreateUser(): ", req)
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}

	err := us.userService.CreateUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create user. {}", err)
	}

	c.Status(http.StatusOK)
}
