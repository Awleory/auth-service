package service

import (
	"context"

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
	Create(ctx context.Context, token domain.RefreshSession) error
	Get(ctx context.Context, token string) (domain.RefreshSession, error)
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
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	user := domain.User{
		Email:    inp.Email,
		Password: password,
	}

	return s.repo.Create(ctx, user)
}

func (s *Users) SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error) {
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.Get(ctx, inp.Email, password)
	if err != nil {
		return "", "", err
	}

	return s.generateTokens(ctx, int64(user.ID))
}
