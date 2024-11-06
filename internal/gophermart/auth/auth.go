package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/storage"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	EXPIREDAT = 30 //время за которое токен истекает в минутах
)

type AuthManager struct {
	s *storage.Storage
	key []byte
}

func New(stor *storage.Storage) (*AuthManager, error) {
	var m AuthManager
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	m.key = key
	m.s = stor

	return &m, nil
}

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (m *AuthManager) Registration(username, password string) (string, error) {
	var user storage.User
	hashPw, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	user.Username = username
	user.Password_hash = hashPw
	user.Created = time.Now()
	err = m.s.NewUser(user)
	if err != nil {
		return "", err
	}
	token, err := m.createToken(username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *AuthManager) Login(username, password string) (string, error) {
	user, err := m.s.GetUser(username)
	if err != nil {
		return "", err
	}
	if !checkPasswordHash(password, user.Password_hash) {
		return "", ErrWrongPassword
	}
	token, err := m.createToken(username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *AuthManager) Authentication(token string) error {
	claims, err := m.decodeToken(token)
	if err != nil {
		return err
	}
	_, err = m.s.GetUser(claims.Username)
	return err
}

func (m *AuthManager) createToken(username, password string) (string, error) {
	expirationTime := time.Now().Add(EXPIREDAT * time.Minute)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenS, err := token.SignedString(m.key)
	if err != nil {
		return "", err
	}

	return tokenS, nil
}

func (m *AuthManager) decodeToken(tokenS string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenS, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
        }
		return m.key, nil
	})

	if err != nil {
        return nil, err
    }

	if !token.Valid {
        return nil, ErrInvalidToken
    }

	return claims, nil
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}