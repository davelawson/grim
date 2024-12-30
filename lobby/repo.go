package lobby

import (
	"database/sql"

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
	result, err := tx.Exec("update lobbies set name = ?, owner_id = ? where id = ?", name, ownerId, lobbyId)
	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), err
}

func (repo *LobbyRepo) AddUserToLobby(tx *sql.Tx, lobbyId string, userId string) error {
	_, err := tx.Exec("insert into lobby_users(lobby_id, user_id) values(?, ?)", lobbyId, userId)
	return err
}

func (repo *LobbyRepo) GetLobby(tx *sql.Tx, id string) (*Lobby, error) {
	row := tx.QueryRow("select id, name, owner_id from lobbies where id = ?", id)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) GetLobbyMembers(tx *sql.Tx, id string) ([]string, error) {
	queryRows, err := tx.Query("select user_id from lobby_users where lobby_id = ?", id)
	if err != nil {
		return nil, err
	}
	var userIds []string
	for queryRows.Next() {
		var userId string
		scanErr := queryRows.Scan(&userId)
		if scanErr != nil {
			return nil, scanErr
		}
		userIds = append(userIds, userId)
	}
	return userIds, nil
}

func (repo *LobbyRepo) GetLobbyByNameAndOwner(name string, ownerId string) (*Lobby, error) {
	row := repo.db.QueryRow("select uuid, name, owner from lobbies where name = ? and owner_id = ?", name, ownerId)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) DeleteLobby(tx *sql.Tx, lobbyId string, ownerId string) (int, error) {
	result, err := tx.Exec("delete from lobbies where id = ? and owner_id = ?", lobbyId, ownerId)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), err
}

func (repo *LobbyRepo) RemoveMemberFromLobby(tx *sql.Tx, lobbyId string, userId string) error {
	_, err := tx.Exec("delete from lobby_users where lobby_id = ? and user_id = ?", lobbyId, userId)
	return err
}

func (repo *LobbyRepo) scanLobby(row *sql.Row) (*Lobby, error) {
	lobby := Lobby{}
	err := row.Scan(&lobby.Id, &lobby.Name, &lobby.Owner)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &lobby, err
}
