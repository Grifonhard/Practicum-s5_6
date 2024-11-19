package model

import (
	"time"
)

// TODO раскидать ордеры по слоям?

// возможные статусы заказов
const (
	NEW        = "NEW"
	PROCESSING = "PROCESSING"
	INVALID    = "INVALID"
	PROCESSED  = "PROCESSED"
)

// возможные статусы заказов в бд
const (
	NEWINT = iota
	PROCESSINGINT
	INVALIDINT
	PROCESSEDINT
)

type User struct {
	ID            int
	Username      string
	PasswordHash string
	Created       time.Time
}

type Order struct {
	ID      int
	UserID  int
	Status  string
	Created time.Time
	Updated time.Time
}

type BalanceTransactions struct {
	ID      int
	UserID  int
	OrderID int
	Sum     int
	Created time.Time
}
