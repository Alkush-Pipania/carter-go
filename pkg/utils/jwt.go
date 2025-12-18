package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Jwtservice struct {
	secret []byte
}

func NewJwtservice(secret string) *Jwtservice {
	return &Jwtservice{
		secret: []byte(secret),
	}
}

func (s *Jwtservice) Verify(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("user_id missing in token")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return "", errors.New("token expired")
		}
	}

	return userID, nil
}
