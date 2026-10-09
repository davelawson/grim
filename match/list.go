package match

import (
	"database/sql"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *Repo) listForUser(tx *sql.Tx, actorID string) ([]api.MatchSummary, error) {
	rows, err := tx.Query("select items.id, items.status from matches items join players membership on membership.match_id = items.id where membership.player_id = ? and items.deleted_at is null order by items.id", actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []api.MatchSummary{}
	for rows.Next() {
		var item api.MatchSummary
		if err := rows.Scan(&item.ID, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) List(tx *sql.Tx, actorID string) (*api.ListMatchesResponse, error) {
	items, err := s.repo.listForUser(tx, actorID)
	if err != nil {
		return nil, err
	}
	return &api.ListMatchesResponse{Matches: items}, nil
}

func (sf *ServiceFacade) List(actorID string) (*api.ListMatchesResponse, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*api.ListMatchesResponse, error) { return sf.service.List(tx, actorID) })()
}

// List returns only IDs and statuses for the authenticated user's memberships.
// @Summary List my matches
// @Security ApiKeyAuth
// @Tags match
// @Produce json
// @Success 200 {object} api.ListMatchesResponse
// @Failure 401,500 {object} api.ErrorResponse
// @Router /matches [get]
func (c *Controller) List(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	response, err := c.service.List(actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
