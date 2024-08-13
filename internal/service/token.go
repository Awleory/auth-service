package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/awleory/medodstest/internal/domain"
	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

const (
	tokenExpires        = 720 // hour
	refreshTokenExpires = 15  // minute
)

type UserClaims struct {
	jwt.MapClaims
	IP  string `json:"ip-address"`
	IP1 string `json:"issuedAt"`
	IP2 string `json:"expiresAt"`
	IP3 string `json:"userID"`
}

func (user *Users) generateTokens(ctx context.Context, userClaims domain.JWTUserClaims) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID":     strconv.Itoa(int(userClaims.ID)),
		"issuedAt":   time.Now().Unix(),
		"expiresAt":  time.Now().Add(time.Minute * refreshTokenExpires).Unix(),
		"ip-address": userClaims.IP,
	})

	accessToken, err := t.SignedString(user.hmacSecret)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := user.sessionsRepo.Create(ctx, domain.RefreshToken{
		UserID:    userClaims.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * tokenExpires),
		UserIP:    userClaims.IP,
	}); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func newRefreshToken() (string, error) {
	b := make([]byte, 32)

	s := rand.NewSource(time.Now().Unix())
	r := rand.New(s)

	if _, err := r.Read(b); err != nil {
		return "", err
	}

	return fmt.Sprintf("%x", b), nil
}

func (user *Users) ParseToken(ctx context.Context, token string) (int64, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return user.hmacSecret, nil
	}

	t, err := jwt.Parse(token, keyFunc)
	if err != nil {
		return 0, err
	}

	if !t.Valid {
		return 0, errors.New("invalid token")
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("invalid claims")
	}

	subject, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid subject")
	}

	id, err := strconv.Atoi(subject)
	if err != nil {
		return 0, errors.New("invalid subject")
	}

	return int64(id), nil
}

func (user *Users) RefreshTokens(ctx context.Context, refreshToken string, userIP string) (string, string, error) {
	session, err := user.sessionsRepo.Get(ctx, refreshToken)
	if err != nil {
		return "", "", err
	}

	if session.ExpiresAt.Unix() < time.Now().Unix() {
		return "", "", errors.New("refresh token expired")
	}

	if session.UserIP != userIP {
		email, err := user.repo.GetEmail(ctx, session.UserID)

		if err == nil {
			SendMsg(email)
		}
		return "", "", errors.New("ip addresses are different")
	}

	return user.generateTokens(ctx, domain.JWTUserClaims{
		ID: session.UserID,
		IP: userIP,
	})
}

func SendMsg(email string) {
	// отправка письма о попытке входа с другого устройства
	logrus.Infof("sending message to %s", email)
}
