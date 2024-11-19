package model

import (
	"encoding/json"
	"strconv"
	"time"
)

type OrderDto struct {
	ID         string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"-"`
	UploadedAt string  `json:"uploaded_at"`
}

func (of *OrderDto) ConvertOrder(o *Order, acc float64) error {
	of.ID = strconv.Itoa(o.ID)
	of.Status = o.Status
	of.Accrual = acc
	of.UploadedAt = o.Updated.Format(time.RFC3339)
	return nil
}

type BalanceDto struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

type WithdrawlDto struct {
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at"`
}

func GetWithdrawFront(order int, sum float64, processed time.Time) *WithdrawlDto {
	var w WithdrawlDto
	w.Order = strconv.Itoa(order)
	w.Sum = sum
	w.ProcessedAt = processed.Format(time.RFC3339)
	return &w
}

func (of OrderDto) MarshalJSON() ([]byte, error) {
	type Alias OrderDto
	aux := struct {
		Alias
		Accrual *float64 `json:"accrual,omitempty"`
	}{
		Alias:   Alias(of),
		Accrual: nil,
	}

	if of.Accrual != 0 {
		aux.Accrual = &of.Accrual
	}

	return json.Marshal(aux)
}

// UnmarshalJSON - кастомное анмаршалирование
func (of *OrderDto) UnmarshalJSON(data []byte) error {
	type Alias OrderDto
	aux := struct {
		Alias
		Accrual *float64 `json:"accrual"`
	}{
		Alias: Alias(*of),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Accrual != nil {
		of.Accrual = *aux.Accrual
	} else {
		of.Accrual = 0
	}

	return nil
}

type WithdrawRequest struct {
	Order string  `json:"order" binding:"required"`
	Sum   float64 `json:"sum" binding:"required"`
}
