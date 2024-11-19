package model

import "time"

type OrderDB struct {
	ID      int
	UserID  int
	Status  int
	Created time.Time
	Updated time.Time
}

func (o *Order) HydrateDB(odb *OrderDB) {
	o.ID = odb.ID
	o.UserID = odb.UserID
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
