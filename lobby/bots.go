package lobby

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/davelawson/grim/api"
	"github.com/davelawson/grim/model"
	"github.com/davelawson/grim/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *LobbyRepo) GetBots(tx *sql.Tx, id string) (map[string]api.Bot, error) {
	rows, err := tx.Query("select id, name from lobby_bots where lobby_id = ? order by id", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	bots := map[string]api.Bot{}
	for rows.Next() {
		var id string
		var bot api.Bot
		if err := rows.Scan(&id, &bot.Name); err != nil {
			return nil, err
		}
		bots[id] = bot
	}
	return bots, rows.Err()
}

// ValidateStorage checks installation without changing an existing database.
func (r *LobbyRepo) ValidateStorage(tx *sql.Tx) error {
	rows, err := tx.Query("select id, lobby_id, name from lobby_bots limit 0")
	if err != nil {
		return err
	}
	return rows.Close()
}

// changeBot shares the same immediate transaction as human roster changes.
func (s *Service) changeBot(tx *sql.Tx, id, actor, removeID string) (*api.Lobby, error) {
	lobby, err := s.GetLobby(tx, id)
	if err != nil {
		return nil, err
	}
	if lobby.Owner != actor {
		return nil, util.ErrForbidden
	}
	if lobby.Status != "open" {
		return nil, util.ErrLobbyClosed
	}
	if removeID != "" {
		if _, ok := lobby.Bots[removeID]; !ok {
			return nil, util.ErrNotFound
		}
		_, err = tx.Exec("delete from lobby_bots where lobby_id = ? and id = ?", id, removeID)
	} else {
		if len(lobby.Members)+len(lobby.Bots) >= maxLobbyMembers {
			return nil, errLobbyFull
		}
		used := map[string]bool{}
		for _, bot := range lobby.Bots {
			used[bot.Name] = true
		}
		n := 1
		for used[fmt.Sprintf("Bot %d", n)] {
			n++
		}
		_, err = tx.Exec("insert into lobby_bots(id, lobby_id, name) values(?, ?, ?)", "bot-"+uuid.NewString(), id, fmt.Sprintf("Bot %d", n))
	}
	if err != nil {
		return nil, err
	}
	if err := s.repo.ClearReadiness(tx, id); err != nil {
		return nil, err
	}
	return s.GetLobby(tx, id)
}

func (sf *ServiceFacade) changeBot(id, actor, removeID string) (*api.Lobby, error) {
	return util.InTxWithResult(sf.db, func(tx *sql.Tx) (*api.Lobby, error) {
		return sf.service.changeBot(tx, id, actor, removeID)
	})()
}

// AddBot adds one server-controlled opponent to an open lobby.
// @Summary Add an automated opponent (owner only)
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Param id path string true "Lobby ID"
// @Success 201 {object} api.GetLobbyResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/bot [post]
func (c *Controller) AddBot(ctx *gin.Context) { c.changeBot(ctx, "", http.StatusCreated) }

// RemoveBot removes a server-controlled opponent from an open lobby.
// @Summary Remove an automated opponent (owner only)
// @Security ApiKeyAuth
// @Tags lobby
// @Produce json
// @Param id path string true "Lobby ID"
// @Param botid path string true "Automated opponent ID"
// @Success 200 {object} api.GetLobbyResponse
// @Failure 401,403,404,409,500 {object} api.ErrorResponse
// @Router /lobby/{id}/bot/{botid} [delete]
func (c *Controller) RemoveBot(ctx *gin.Context) { c.changeBot(ctx, ctx.Param("botid"), http.StatusOK) }

func (c *Controller) changeBot(ctx *gin.Context, removeID string, status int) {
	actor := ctx.MustGet("reqUser").(*model.User)
	view, err := c.lobbyService.changeBot(ctx.Param("id"), actor.Id, removeID)
	if err != nil {
		util.WriteAPIError(ctx, err)
		return
	}
	ctx.JSON(status, api.GetLobbyResponse{Lobby: *view})
}
