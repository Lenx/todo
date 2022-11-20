package repository

import (
	"fmt"

	"github.com/Hanqur/todo_app"
	"github.com/jmoiron/sqlx"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// отправляем запрос на добавление нового пользователя в БД и получаем его Id
func (r *AuthPostgres) CreateUser(user todo_app.User) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", userTable)

	row := r.db.QueryRow(query, user.Name, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

// отправляем запрос на получения пользователя из БД
func (r *AuthPostgres) GetUser(username, password string) (todo_app.User, error) {
	var user todo_app.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", userTable)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		fmt.Println(err)
	}

	return user, err
}