package user

import (
	"database/sql"
)

type Repo struct {
	db *sql.DB
}

func NewRepo(db *sql.DB) *Repo {
	return &Repo{db: db}
}

func (repo *Repo) CreateUser(email string, name string, passwordHash []byte) error {
	_, err := repo.db.Exec(
		"insert into users(email, name, password_hash) values(?, ?, ?)",
		email, name, passwordHash)
	return err
}

func (repo *Repo) GetUserByEmail(email string) (*User, error) {
	row := repo.db.QueryRow("select id, email, name, password_hash from users where email = ?", email)
	return repo.scanUser(row)
}

func (repo *Repo) scanUser(row *sql.Row) (*User, error) {
	user := User{}
	err := row.Scan(&user.Id, &user.Email, &user.Name, &user.PasswordHash)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &user, err
}
