package model

import (
	"strconv"
	"time"
)

type OrderDto struct {
	Id         string `json:"number"`
	Status     string `json:"status"`
	Accrual    int    `json:"accrual"`
	UploadedAt string `json:"uploaded_at"`
}

func (of *OrderDto) ConvertOrder(o *Order, acc int) error {
	of.Id = strconv.Itoa(o.Id)
	of.Status = o.Status
	of.Accrual = acc
	of.UploadedAt = o.Updated.Format(time.RFC3339)
	return nil
}

type BalanceDto struct {
	Current   int `json:"current"`
	Withdrawn int `json:"withdrawn"`
}

type WithdrawlDto struct {
	Order       string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at"`
}

func GetWithdrawFront(order, sum int, processed time.Time) *WithdrawlDto {
	var w WithdrawlDto
	w.Order = strconv.Itoa(order)
	w.Sum = sum
	w.ProcessedAt = processed.Format(time.RFC3339)
	return &w
}
