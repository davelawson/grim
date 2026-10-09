package lobby

import (
	"database/sql"
	"main/util"

	"github.com/google/uuid"
)

type LobbyRepo struct {
	db *sql.DB
}

func NewLobbyRepo(db *sql.DB) *LobbyRepo {
	return &LobbyRepo{db: db}
}

func (repo *LobbyRepo) CreateLobby(tx *sql.Tx, name string, ownerId string) (string, error) {
	newUuid := uuid.New()
	_, err := tx.Exec("insert into lobbies(id, name, owner_id) values(?, ?, ?)", newUuid.String(), name, ownerId)
	if err != nil {
		return "", err
	}
	return newUuid.String(), nil
}

func (repo *LobbyRepo) UpdateLobby(tx *sql.Tx, lobbyId string, name string, ownerId string) (int, error) {
	result, err := tx.Exec("update lobbies set name = ?, owner_id = ? where id = ? and deleted_at is null", name, ownerId, lobbyId)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	return int(rowsAffected), err
}

func (repo *LobbyRepo) AddUserToLobby(tx *sql.Tx, lobbyId string, userId string) error {
	_, err := tx.Exec("insert into lobby_users(lobby_id, user_id) values(?, ?)", lobbyId, userId)
	return err
}

func (repo *LobbyRepo) GetLobby(tx *sql.Tx, id string) (*Lobby, error) {
	row := tx.QueryRow("select id, name, owner_id, status, match_id from lobbies where id = ? and deleted_at is null", id)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) GetLobbyMembers(tx *sql.Tx, id string) ([]string, error) {
	queryRows, err := tx.Query("select user_id from lobby_users where lobby_id = ? order by user_id", id)
	if err != nil {
		return nil, err
	}
	defer queryRows.Close()
	userIds := []string{}
	for queryRows.Next() {
		var userId string
		scanErr := queryRows.Scan(&userId)
		if scanErr != nil {
			return nil, scanErr
		}
		userIds = append(userIds, userId)
	}
	return userIds, queryRows.Err()
}

func (repo *LobbyRepo) GetLobbyByNameAndOwner(name string, ownerId string) (*Lobby, error) {
	row := repo.db.QueryRow("select id, name, owner_id, status, match_id from lobbies where name = ? and owner_id = ? and deleted_at is null", name, ownerId)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) DeleteLobby(tx *sql.Tx, lobbyId string, ownerId string) (int, error) {
	result, err := tx.Exec("update lobbies set deleted_at = datetime('now') where id = ? and owner_id = ? and deleted_at is null", lobbyId, ownerId)
	if err != nil {
		return 0, err
	}
	rowsAffected, err := result.RowsAffected()
	return int(rowsAffected), err
}

func (repo *LobbyRepo) RemoveMemberFromLobby(tx *sql.Tx, lobbyId string, userId string) error {
	_, err := tx.Exec("delete from lobby_users where lobby_id = ? and user_id = ?", lobbyId, userId)
	return err
}

func (repo *LobbyRepo) scanLobby(row *sql.Row) (*Lobby, error) {
	lobby := Lobby{}
	err := row.Scan(&lobby.Id, &lobby.Name, &lobby.Owner, &lobby.Status, &lobby.MatchID)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &lobby, err
}

func (repo *LobbyRepo) GetReadiness(tx *sql.Tx, id string) (map[string]bool, error) {
	rows, err := tx.Query("select user_id, ready from lobby_users where lobby_id = ?", id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	readiness := make(map[string]bool)
	for rows.Next() {
		var userID string
		var ready bool
		if err := rows.Scan(&userID, &ready); err != nil {
			return nil, err
		}
		readiness[userID] = ready
	}
	return readiness, rows.Err()
}

func (repo *LobbyRepo) ClearReadiness(tx *sql.Tx, id string) error {
	_, err := tx.Exec("update lobby_users set ready = 0 where lobby_id = ?", id)
	return err
}

func (repo *LobbyRepo) SetReady(tx *sql.Tx, id, userID string, ready bool) error {
	_, err := tx.Exec("update lobby_users set ready = ? where lobby_id = ? and user_id = ?", ready, id, userID)
	return err
}

func (repo *LobbyRepo) Close(tx *sql.Tx, id, matchID string) error {
	result, err := tx.Exec("update lobbies set status = 'closed', match_id = ? where id = ? and status = 'open' and deleted_at is null", matchID, id)
	if err != nil {
		return err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return util.ErrLobbyClosed
	}
	return nil
}
