package user

import (
	"fmt"
	"main/model"
	"main/model/api"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type authService interface {
	VerifyBearerToken(token string) (*model.User, error)
}

type Controller struct {
	userService *ServiceFacade
}

func NewController(userServiceFacade *ServiceFacade) *Controller {
	return &Controller{
		userService: userServiceFacade,
	}
}

// @Summary		Get user by email
// @Description	Lookup a specific user by email
// @Security ApiKeyAuth
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			request	body		api.GetUserRequest	true	"Request Object"
// @Success		200		{object}	api.GetUserResponse
// @Router			/user/getbyemail [post]
func (us *Controller) GetUserByEmail(c *gin.Context) {
	// TODO: find a re-usable way to translate the context into a typed request
	req := api.GetUserRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	fmt.Println("GetUserByEmail(): ", req)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Unable to interpret payload", reqErr)
		return
	}
	req.Email = strings.ToLower(req.Email)

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

// @Summary		Create user
// @Description	Create a new user
// @Tags			user
// @Accept			json
// @Produce		json
// @Param			request	body	api.CreateUserRequest	true	"Request Object"
// @Success		200
// @Router			/user [post]
func (us *Controller) CreateUser(c *gin.Context) {
	req := api.CreateUserRequest{}
	reqErr := c.ShouldBindBodyWith(&req, binding.JSON)
	if reqErr != nil {
		c.String(http.StatusBadRequest, "Invalid request body")
		c.Error(reqErr)
		return
	}
	req.Email = strings.ToLower(req.Email)
	fmt.Println("CreateUser(): ", req)

	err := us.userService.CreateUser(req.Email, req.Name, req.Password)
	if err != nil {
		c.String(http.StatusInternalServerError, "Something went wrong.  Unable to create user. {}", err)
	}

	c.Status(http.StatusOK)
}
