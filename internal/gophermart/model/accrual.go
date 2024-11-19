package model

import "strconv"

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
	status, err := strconv.Atoi(oA.Status)
	return oA.Accrual, status, err
}
