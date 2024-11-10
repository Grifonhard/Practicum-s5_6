package model

import "time"

type User struct {
	Id            int
	Username      string
	Password_hash string
	Created       time.Time
}

type OrderDB struct {
	Id      int
	UserId  int
	OrderId int
	Status  int
	Created time.Time
}

type Order struct {
	Id      int
	UserId  int
	OrderId int
	Status  string
	Created time.Time
}

func (o *Order) Convert(odb *OrderDB) {
	o.Id = odb.Id
	o.UserId = odb.UserId
	o.OrderId = odb.OrderId
	o.Created = odb.Created
	switch odb.Status {
	case 0:
		o.Status = NEW
	case 1:
		o.Status = PROCESSING
	case 2:
		o.Status = INVALID
	case 3:
		o.Status = PROCESSED
	default:
		// TODO запись в логи
		o.Status = INVALID
	}
}

type BalTrans struct {
	Id      int
	UserId  int
	OrderId int
	Sum     int
	Created time.Time
}

// возможные статусы заказов
const (
	NEW = "NEW"
	PROCESSING = "PROCESSING"
	INVALID = "INVALID"
	PROCESSED = "PROCESSED"
)