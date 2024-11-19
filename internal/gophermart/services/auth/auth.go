package auth

import (
	"crypto/rand"
	"fmt"
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/repository"
	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/services/transactions"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	EXPIREDAT = 30 //время за которое токен истекает в минутах
)

type Manager struct {
	p         *repository.DB
	muT		  *transactions.Mutex
	secretKey []byte
}

func New(db *repository.DB, t *transactions.Mutex) (*Manager, error) {
	var m Manager
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}

	m.muT = t
	m.secretKey = key
	m.p = db

	return &m, nil
}

// TODO TokenClaims
type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func (m *Manager) Registration(username, password string) (string, error) {

	m.muT.Lock(username)
	defer m.muT.Unlock(username)

	hashPw, err := hashPassword(password)

	if err != nil {
		return "", err
	}
	err = m.p.InsertUser(username, hashPw)
	if err != nil {
		return "", err
	}
	user, err := m.p.GetUser(username)
	if err != nil {
		return "", err
	}
	token, err := m.createToken(user.ID)
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
	if !checkPasswordHash(password, user.PasswordHash) {
		return "", ErrWrongPassword
	}
	token, err := m.createToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *Manager) Authentication(token string) (int, error) {

	claims, err := m.decodeToken(token)

	if err != nil {
		return 0, err
	}
	user, err := m.p.GetUserByID(claims.UserID)
	return user.ID, err
}

func (m *Manager) createToken(userID int) (string, error) {
	expirationTime := time.Now().Add(EXPIREDAT * time.Minute)

	claims := &Claims{
		UserID: userID,
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
