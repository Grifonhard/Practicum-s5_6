package model

import (
	"strconv"
	"time"
)

type User struct {
	Id            int
	Username      string
	Password_hash string
	Created       time.Time
}

type OrderDB struct {
	Id      int
	UserId  int
	Status  int
	Created time.Time
	Updated time.Time
}

type Order struct {
	Id      int
	UserId  int
	Status  string
	Created time.Time
	Updated time.Time
}

func (o *Order) HydrateDB(odb *OrderDB) {
	o.Id = odb.Id
	o.UserId = odb.UserId
	o.Created = odb.Created
	switch odb.Status {
	case NEWINT:
		o.Status = NEW
	case PROCESSINGINT:
		o.Status = PROCESSING
	case INVALIDINT:
		o.Status = INVALID
	case PROCESSEDINT:
		o.Status = PROCESSED
	default:
		// TODO запись в логи
		o.Status = INVALID
	}
}

type OrderAccrual struct {
	OrderId string `json:"order"`
	Status string `json:"status"`
	Accrual int `json:"accrual"`
}

func (o *Order) ConvertAccrual(oA *OrderAccrual) (int, int, error) {
	var err error
	o.Id, err = strconv.Atoi(oA.OrderId)
	if err != nil {
		return 0, 0, err
	}
	o.Status = oA.Status
	status, err := strconv.Atoi(oA.Status)
	return oA.Accrual, status, err
}

// TODO раскидать ордеры по слоям
type OrderDto struct {
	Id string `json:"number"`
	Status string `json:"status"`
	Accrual int `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

func (of *OrderDto) ConvertOrder(o *Order, acc int) error {
	of.Id = strconv.Itoa(o.Id)
	of.Status = o.Status
	of.Accrual = acc
	of.UploadedAt = o.Updated.Format(time.RFC3339)
	return nil
}

type BalanceTransactions struct {
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

const (
	NEWINT = iota
	PROCESSINGINT
	INVALIDINT
	PROCESSEDINT
)