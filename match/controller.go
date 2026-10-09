package match

import (
	"main/game"
	"main/model"
	"main/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Match *game.View `json:"match"`
}

type Controller struct{ service *ServiceFacade }

func NewController(service *ServiceFacade) *Controller { return &Controller{service: service} }

// Get returns the authenticated participant's lifecycle view.
// @Summary Get a match
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} match.Response
// @Failure 401,403,404,409,500 {object} util.ErrorResponse
// @Router /match/{id} [get]
func (c *Controller) Get(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.service.Get(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, Response{Match: view})
}

// End soft-deletes a match while preserving its saved state.
// @Summary Administratively end a match
// @Description Requires current stored administrator permission; already-deleted matches succeed unchanged.
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param id path string true "Match ID"
// @Success 204 "No Content"
// @Failure 401,403,404,500 {object} util.ErrorResponse
// @Router /admin/match/{id}/end [post]
func (c *Controller) End(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	if err := c.service.End(ctx.Param("id"), actor.Id); err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
