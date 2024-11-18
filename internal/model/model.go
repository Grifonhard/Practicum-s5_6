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
	Id            int
	Username      string
	Password_hash string
	Created       time.Time
}

type Order struct {
	Id      int
	UserId  int
	Status  string
	Created time.Time
	Updated time.Time
}

type BalanceTransactions struct {
	Id      int
	UserId  int
	OrderId int
	Sum     int
	Created time.Time
}
