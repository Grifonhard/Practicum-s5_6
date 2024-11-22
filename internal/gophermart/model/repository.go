package model

import (
	"time"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
)

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
		logger.Error("invalid status when convert db: %d", odb.Status)
		o.Status = INVALID
	}
}
