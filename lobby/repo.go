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

func (repo *LobbyRepo) CreateLobby(name string, ownerId string) (string, error) {
	newUuid := uuid.New()
	_, err := repo.db.Exec("insert into lobbies(id, name, owner_id) values(?, ?, ?)", newUuid.String(), name, ownerId)
	if err != nil {
		return "", err
	}
	return newUuid.String(), nil
}

func (repo *LobbyRepo) GetLobby(id string) (*Lobby, error) {
	row := repo.db.QueryRow("select id, name, owner_id from lobbies where id = ?", id)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) GetLobbyByNameAndOwner(name string, ownerId string) (*Lobby, error) {
	row := repo.db.QueryRow("select uuid, name, owner from lobbies where name = ? and owner_id = ?", name, ownerId)
	return repo.scanLobby(row)
}

func (repo *LobbyRepo) DeleteLobby(lobbyId string, ownerId string) (int, error) {
	result, err := repo.db.Exec("delete from lobbies where id = ? and owner_id = ?", lobbyId, ownerId)
	if err != nil {
		return 0, err
	}
	rowsAffected, _ := result.RowsAffected()
	return int(rowsAffected), err
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
