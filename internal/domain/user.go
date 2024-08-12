package domain

import (
	"github.com/go-playground/validator"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

type JWTUserClaims struct {
	ID int64
	IP string
}

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=6"`
}

func (s *SignUpInput) Validate() error {
	return validate.Struct(s)
}

type SignInInput struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,gte=8"`
	IP       string `json:"ip"`
}

func (s *SignInInput) Validate() error {
	return validate.Struct(s)
}
