package psql

import (
	"context"
	"database/sql"

	"github.com/awleory/medodstest/internal/domain"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{db}
}

func (r *Users) Create(ctx context.Context, user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, password_hash) values ($1, $2)",
		user.Email, user.Password)

	return err
}

func (r *Users) Get(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, email FROM users WHERE email=$1 AND password_hash=$2", email, password).
		Scan(&user.ID, &user.Email)

	return user, err
}
