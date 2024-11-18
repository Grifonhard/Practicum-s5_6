package model

import "strconv"

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
