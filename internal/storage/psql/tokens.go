package psql

import (
	"context"
	"database/sql"

	"github.com/awleory/medodstest/internal/domain"
)

type Tokens struct {
	db *sql.DB
}

func NewTokens(db *sql.DB) *Tokens {
	return &Tokens{
		db: db,
	}
}

func (r *Tokens) Create(ctx context.Context, token domain.RefreshToken) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", token.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("INSERT INTO refresh_tokens (user_id, token, expires_at, user_ip) values ($1, $2, $3, $4)",
		token.UserID, token.Token, token.ExpiresAt, token.UserIP)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *Tokens) Get(ctx context.Context, token string) (domain.RefreshToken, error) {
	var t domain.RefreshToken
	tx, err := r.db.Begin()
	if err != nil {
		return t, err
	}

	err = tx.QueryRow("SELECT id, user_id, token, expires_at, user_ip FROM refresh_tokens WHERE token=$1", token).
		Scan(&t.ID, &t.UserID, &t.Token, &t.ExpiresAt, &t.UserIP)
	if err != nil {
		tx.Rollback()
		return t, err
	}

	_, err = tx.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", t.UserID)
	if err != nil {
		tx.Rollback()
		return t, err
	}
	err = tx.Commit()
	return t, err
}
