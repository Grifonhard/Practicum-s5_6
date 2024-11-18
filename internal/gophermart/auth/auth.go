package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/order/storage"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	EXPIREDAT = 30 //время за которое токен истекает в минутах
)

type Manager struct {
	s         *storage.Storage
	p         *repository.DB
	secretKey []byte
}

func New(db *repository.DB, stor *storage.Storage) (*Manager, error) {
	var m Manager
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	m.secretKey = key
	m.s = stor
	m.p = db

	return &m, nil
}

// TODO TokenClaims
type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func (m *Manager) Registration(username, password string) (string, error) {
	hashPw, err := hashPassword(password)
	if err != nil {
		return "", err
	}
	err = m.p.InsertUser(username, hashPw)
	if err != nil {
		return "", err
	}
	token, err := m.createToken(username, password)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *Manager) Login(username, password string) (string, error) {
	user, err := m.p.GetUser(username)
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

func (m *Manager) Authentication(token string) error {
	claims, err := m.decodeToken(token)
	if err != nil {
		return err
	}
	_, err = m.p.GetUser(claims.Username)
	return err
}

func (m *Manager) createToken(username, password string) (string, error) {
	expirationTime := time.Now().Add(EXPIREDAT * time.Minute)

	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenS, err := token.SignedString(m.secretKey)
	if err != nil {
		return "", err
	}

	return tokenS, nil
}

func (m *Manager) decodeToken(tokenS string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenS, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secretKey, nil
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
