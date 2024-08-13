package service

import (
	"context"
	"strings"

	"github.com/awleory/medodstest/internal/domain"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	Get(ctx context.Context, email, password string) (domain.User, error)
}

type SessionsRepository interface {
	Create(ctx context.Context, token domain.RefreshToken) error
	Get(ctx context.Context, token string) (domain.RefreshToken, error)
}

type Users struct {
	repo         UsersRepository
	sessionsRepo SessionsRepository
	hasher       PasswordHasher

	hmacSecret []byte
}

func NewUsers(repo UsersRepository, sessionsRepo SessionsRepository, hasher PasswordHasher, secret []byte) *Users {
	return &Users{
		repo:         repo,
		sessionsRepo: sessionsRepo,
		hasher:       hasher,
		hmacSecret:   secret,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	password_hash, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    strings.ToLower(inp.Email),
		Password: password_hash,
	}

	return s.repo.Create(ctx, user)
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password_hash, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.Get(ctx, strings.ToLower(inp.Email), password_hash)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, domain.JWTUserClaims{ID: int64(user.ID), IP: inp.IP})
}
