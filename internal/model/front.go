package model

import (
	"strconv"
	"time"
)

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