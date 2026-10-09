package match

import (
	"database/sql"
	"encoding/json"
	"errors"
	"main/game"
	"main/util"
)

type Repo struct{ db *sql.DB }

func NewRepo(db *sql.DB) *Repo { return &Repo{db: db} }

type storedMatch struct {
	id, name, lobbyID, status string
	revision                  int64
	outcome                   sql.NullString
	snapshot                  []byte
}

func (r *Repo) load(tx *sql.Tx, id string) (*storedMatch, error) {
	var stored storedMatch
	err := tx.QueryRow(`select id, name, lobby_id, status, revision, outcome, snapshot
		from matches where id = ? and deleted_at is null`, id).Scan(&stored.id, &stored.name, &stored.lobbyID,
		&stored.status, &stored.revision, &stored.outcome, &stored.snapshot)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, util.ErrNotFound
	}
	return &stored, err
}

func (r *Repo) isParticipant(tx *sql.Tx, id, actorID string) (bool, error) {
	var member bool
	err := tx.QueryRow(`select exists(select 1 from players where match_id = ? and player_id = ?)`, id, actorID).Scan(&member)
	return member, err
}

func (r *Repo) create(tx *sql.Tx, state *game.State) error {
	snapshot, err := json.Marshal(state)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`insert into matches(id, name, lobby_id, status, revision, outcome, snapshot)
		values(?, ?, ?, ?, ?, ?, ?)`, state.ID, state.Name, state.LobbyID,
		state.Status, state.Revision, state.Outcome, string(snapshot))
	if err != nil {
		return err
	}
	for _, player := range state.Participants {
		if _, err := tx.Exec("insert into players(player_id, match_id) values(?, ?)", player, state.ID); err != nil {
			return err
		}
	}
	return nil
}

func (r *Repo) save(tx *sql.Tx, state *game.State, expectedRevision int64) error {
	snapshot, err := json.Marshal(state)
	if err != nil {
		return err
	}
	result, err := tx.Exec(`update matches set status = ?, revision = ?, outcome = ?, snapshot = ?
		where id = ? and revision = ? and deleted_at is null`, state.Status, state.Revision, state.Outcome,
		string(snapshot), state.ID, expectedRevision)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rows != 1 {
		return errors.New("match revision changed during transaction")
	}
	return nil
}

func (r *Repo) launchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string) (*game.View, error) {
	var command string
	var data []byte
	var deletedAt sql.NullString
	err := tx.QueryRow(`select receipt.command, receipt.response, matches.deleted_at
		from launch_receipts as receipt join matches on matches.id = receipt.match_id
		where receipt.lobby_id = ? and receipt.actor_id = ? and receipt.request_id = ?`,
		lobbyID, actorID, requestID).Scan(&command, &data, &deletedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if deletedAt.Valid {
		return nil, util.ErrNotFound
	}
	if command != "launch" {
		return nil, &util.APIError{Status: 409, Code: "request_id_conflict", Message: "Request ID already used for different input"}
	}
	var view game.View
	if err := json.Unmarshal(data, &view); err != nil {
		return nil, err
	}
	return &view, nil
}

func (r *Repo) softDelete(tx *sql.Tx, id string) error {
	result, err := tx.Exec("update matches set deleted_at = datetime('now') where id = ? and deleted_at is null", id)
	if err != nil {
		return err
	}
	rows, err := result.RowsAffected()
	if err != nil || rows != 0 {
		return err
	}
	// A previously deleted record is a successful retry; an unknown ID is not.
	var exists bool
	if err := tx.QueryRow("select exists(select 1 from matches where id = ?)", id).Scan(&exists); err != nil {
		return err
	}
	if !exists {
		return util.ErrNotFound
	}
	return nil
}

func (r *Repo) saveLaunchReceipt(tx *sql.Tx, lobbyID, actorID, requestID string, view *game.View) error {
	data, err := json.Marshal(view)
	if err != nil {
		return err
	}
	_, err = tx.Exec(`insert into launch_receipts(lobby_id, actor_id, request_id, match_id, command, response)
		values(?, ?, ?, ?, 'launch', ?)`, lobbyID, actorID, requestID, view.ID, string(data))
	return err
}
