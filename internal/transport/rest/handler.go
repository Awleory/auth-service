package rest

import (
	"context"
	"net/http"

	"github.com/awleory/medodstest/internal/domain"
	"github.com/gorilla/mux"
)

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, accessToken string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type Handler struct {
	usersService User
}

func NewHandler(users User) *Handler {
	return &Handler{
		usersService: users,
	}
}

func (h *Handler) InitRouter() *mux.Router {
	router := mux.NewRouter()
	router.Use(loggingMiddleware)

	auth := router.PathPrefix("/auth").Subrouter()
	{
		auth.HandleFunc("/sign-up", h.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", h.signIn).Methods(http.MethodGet)
		auth.HandleFunc("/refresh", h.refresh).Methods(http.MethodGet)
	}

	return router
}
