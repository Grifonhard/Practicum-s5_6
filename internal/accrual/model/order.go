package model

import "time"

const (
	OrderStatusRegistered = "REGISTERED"
	OrderStatusInvalid    = "INVALID"
	OrderStatusProcessing = "PROCESSING"
	OrderStatusProcessed  = "PROCESSED"
)

type Order struct {
	Number    uint64    `db:"number" json:"number" required:"true"`
	Status    string    `db:"status" json:"status"`
	Accrual   *float64  `db:"accrual" json:"accrual"`
	CreatedAt time.Time `db:"created_at" json:"-"`
}
