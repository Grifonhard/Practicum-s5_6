package model

import (
	"strconv"

	"github.com/Grifonhard/Practicum-s5_6/internal/gophermart/logger"
)

type OrderAccrual struct {
	OrderID string `json:"order"`
	Status  string `json:"status"`
	Accrual int    `json:"accrual"`
}

func (o *Order) ConvertAccrual(oA *OrderAccrual) (int, int, error) {
	var err error
	o.ID, err = strconv.Atoi(oA.OrderID)
	if err != nil {
		return 0, 0, err
	}
	o.Status = oA.Status
	var status int
	switch o.Status {
	case NEW:
		status = NEWINT
	case PROCESSING:
		status = PROCESSINGINT
	case INVALID:
		status = INVALIDINT
	case PROCESSED:
		status = PROCESSEDINT
	default:
		logger.Error("invalid status when convert: %s", o.Status)
		status = 10
	}
	return oA.Accrual, status, nil
}
