package repo

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

func (repo *UserRepo) CreateUser(email string, name string, passwordHash []byte) error {
	_, err := repo.db.Exec(
		"insert into users(id, email, name, password_hash) values(?, ?, ?, ?)",
		uuid.New().String(), email, name, passwordHash)
	return err
}

func (repo *UserRepo) GetUserByEmail(email string) (*model.User, error) {
	row := repo.db.QueryRow("select id, email, name, password_hash from users where email = ?", email)
	return repo.scanUser(row)
}

func (repo *UserRepo) GetUserByToken(token string) (*model.User, error) {
	row := repo.db.QueryRow("select id, email, name, password_hash from users where token = ?", token)
	return repo.scanUser(row)
}

func (repo *UserRepo) UpdateUser(user *model.User) error {
	_, err := repo.db.Exec("update users set name = ?, password_hash = ?, token = ? where id = ?", user.Name, user.PasswordHash, user.Token, user.Id)
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
