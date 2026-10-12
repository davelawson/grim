package lobby

import (
	"database/sql"
	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/util"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (r *LobbyRepo) listForUser(tx *sql.Tx, actorID string) ([]api.LobbySummary, error) {
	rows, err := tx.Query("select items.id, items.status from lobbies items join lobby_users membership on membership.lobby_id = items.id where membership.user_id = ? and items.deleted_at is null order by items.id", actorID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []api.LobbySummary{}
	for rows.Next() {
		var item api.LobbySummary
		if err := rows.Scan(&item.ID, &item.Status); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) List(tx *sql.Tx, actorID string) (*api.ListLobbiesResponse, error) {
	items, err := s.repo.listForUser(tx, actorID)
	if err != nil {
		return nil, err
	}
	return &api.ListLobbiesResponse{Lobbies: items}, nil
}

func (sf *ServiceFacade) List(actorID string) (*api.ListLobbiesResponse, error) {
	return util.InTypedTx(sf.db, func(tx *sql.Tx) (*api.ListLobbiesResponse, error) { return sf.service.List(tx, actorID) })()
}

// List returns only IDs and statuses for the authenticated user's memberships.
// @Summary List my lobbies
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Success 200 {object} api.ListLobbiesResponse
// @Failure 401,500 {object} api.ErrorResponse
// @Router /lobbies [get]
func (c *Controller) List(ctx *gin.Context) {
	actor := ctx.MustGet("reqUser").(*model.User)
	response, err := c.lobbyService.List(actor.Id)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}

func (r *LobbyRepo) listPublic(tx *sql.Tx) ([]api.PublicLobbySummary, error) {
	rows, err := tx.Query(`select lobby.id, lobby.name, lobby.owner_id, count(membership.user_id)
 from lobbies lobby left join lobby_users membership on membership.lobby_id = lobby.id
 where lobby.status = 'open' and lobby.deleted_at is null
 group by lobby.id having count(membership.user_id) < ? order by lobby.id`, maxLobbyMembers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []api.PublicLobbySummary{}
	for rows.Next() {
		var item api.PublicLobbySummary
		if err := rows.Scan(&item.ID, &item.Name, &item.Owner, &item.MemberCount); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Service) ListPublic(tx *sql.Tx) (*api.ListPublicLobbiesResponse, error) {
	items, err := s.repo.listPublic(tx)
	if err != nil {
		return nil, err
	}
	return &api.ListPublicLobbiesResponse{Lobbies: items}, nil
}

func (sf *ServiceFacade) ListPublic() (*api.ListPublicLobbiesResponse, error) {
	return util.InTypedTx(sf.db, sf.service.ListPublic)()
}

// ListPublic returns available public lobby summaries to authenticated users.
// @Summary List public lobbies with free slots
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Success 200 {object} api.ListPublicLobbiesResponse
// @Failure 401,500 {object} api.ErrorResponse
// @Router /lobbies/public [get]
func (c *Controller) ListPublic(ctx *gin.Context) {
	response, err := c.lobbyService.ListPublic()
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, response)
}
