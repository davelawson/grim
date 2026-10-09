package match

import (
	"encoding/json"
	"io"
	"main/api"
	"main/model"
	"main/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Controller struct{ service *ServiceFacade }

func NewController(service *ServiceFacade) *Controller { return &Controller{service: service} }

// Get returns the authenticated participant's private and public match view.
// @Summary Get a match
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} api.MatchResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id} [get]
func (c *Controller) Get(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.service.Get(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.MatchResponse{Match: view})
}

// Command locks one opening draft choice, or replays an identical accepted retry.
// @Summary Choose an opening card
// @Security ApiKeyAuth
// @Tags match
// @Accept json
// @Produce json
// @Param id path string true "Match ID"
// @Param command body api.CommandRequest true "Draft command; expectedrevision is required"
// @Success 200 {object} api.MatchResponse
// @Failure 400,401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/commands [post]
func (c *Controller) Command(ctx *gin.Context) {
	var command api.CommandRequest
	decoder := json.NewDecoder(ctx.Request.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&command); err != nil {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		util.WriteAPIError(ctx, util.ErrInvalidRequest)
		return
	}
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.service.Command(ctx.Param("id"), actor.Id, command)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.MatchResponse{Match: view})
}

// Catalogue returns the definitions pinned to this participant's match.
// @Summary Get the match card catalogue
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Param id path string true "Match ID"
// @Success 200 {object} api.CatalogueResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /match/{id}/catalogue [get]
func (c *Controller) Catalogue(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	catalogue, err := c.service.Catalogue(ctx.Param("id"), actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, api.CatalogueResponse{Catalogue: catalogue})
}

// End soft-deletes a match while preserving its saved state.
// @Summary Administratively end a match
// @Description Requires current stored administrator permission; already-deleted matches succeed unchanged.
// @Security ApiKeyAuth
// @Tags admin
// @Produce json
// @Param id path string true "Match ID"
// @Success 204 "No Content"
// @Failure 401,403,404,500 {object} api.ErrorResponse
// @Router /admin/match/{id}/end [post]
func (c *Controller) End(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	if err := c.service.End(ctx.Param("id"), actor.Id); err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}
