package auth

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	EXPIREDAT = 30 //время за которое токен истекает в минутах
)

var key = []byte("top_secret")

type AuthManager struct {
	
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func createToken(username string) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)

	Claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims)

	tokenS, err := token.SignedString(key)
	if err != nil {
		return "", err
	}

	return tokenS, nil
}

func decodeToken(tokenS string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenS, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
		return key, nil
	})

	if err != nil {
        return nil, err
    }

	if !token.Valid {
        return nil, fmt.Errorf("invalid token")
    }

	return claims, nil
}