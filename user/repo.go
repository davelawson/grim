package user

import (
	"database/sql"
	"main/model"

	uuid "github.com/google/uuid"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func (repo *UserRepo) CreateUser(tx *sql.Tx, email string, name string, passwordHash []byte) error {
	_, err := tx.Exec(
		"insert into users(id, email, name, password_hash) values(?, ?, ?, ?)",
		uuid.New().String(), email, name, passwordHash)
	return err
}

func (repo *UserRepo) GetUserById(tx *sql.Tx, id string) (*model.User, error) {
	row := tx.QueryRow("select id, email, name, password_hash from users where id = ?", id)
	return repo.scanUser(row)
}

func (repo *UserRepo) GetUserByEmail(tx *sql.Tx, email string) (*model.User, error) {
	row := tx.QueryRow("select id, email, name, password_hash from users where email = ?", email)
	return repo.scanUser(row)
}

func (repo *UserRepo) GetUserByToken(tx *sql.Tx, token string) (*model.User, error) {
	row := tx.QueryRow("select id, email, name, password_hash from users where token = ?", token)
	return repo.scanUser(row)
}

func (repo *UserRepo) UpdateUser(tx *sql.Tx, user *model.User) error {
	_, err := tx.Exec("update users set name = ?, password_hash = ?, token = ? where id = ?", user.Name, user.PasswordHash, user.Token, user.Id)
	return err
}

func (repo *UserRepo) scanUser(row *sql.Row) (*model.User, error) {
	user := model.User{}
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, err
}
