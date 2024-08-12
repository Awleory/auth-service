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
)

type UserClaims struct {
	jwt.MapClaims
	IP  string `json:"ip-address"`
	IP1 string `json:"issuedAt"`
	IP2 string `json:"expiresAt"`
	IP3 string `json:"userID"`
}

func (user *Users) generateTokens(ctx context.Context, userClaims domain.JWTUserClaims, compareWithOld bool) (string, string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"userID":     strconv.Itoa(int(userClaims.ID)),
		"issuedAt":   time.Now().Unix(),
		"expiresAt":  time.Now().Add(time.Minute * 15).Unix(),
		"ip-address": userClaims.IP,
	})

	accessToken, err := t.SignedString(user.hmacSecret)
	if err != nil {
		return "", "", err
	}

	if compareWithOld {
		accessTokenOld, err := user.sessionsRepo.GetByID(ctx, userClaims.ID)
		fmt.Println(accessTokenOld, err)
		fmt.Println(userClaims.ID)
		if err == nil {
			if accessTokenOld.Token != accessToken {
				SendMsg()
				return "", "", fmt.Errorf("the IP address are different")
			}
		}
	}

	refreshToken, err := newRefreshToken()
	if err != nil {
		return "", "", err
	}

	if err := user.sessionsRepo.Create(ctx, domain.RefreshSession{
		UserID:    userClaims.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Hour * 24 * 30),
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

	return user.generateTokens(ctx, domain.JWTUserClaims{
		ID: session.UserID,
		IP: userIP,
	}, true)
}

func SendMsg() {
	// отправка письма о попытке входа с другого устройства
}
