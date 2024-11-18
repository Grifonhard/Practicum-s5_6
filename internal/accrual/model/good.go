package model

import (
	"time"
)

type Good struct {
	ID          uint64    `db:"id" json:"-"`
	Description string    `db:"description" json:"description" required:"true"`
	Price       uint64    `db:"price" json:"price" required:"true"`
	OrderNumber uint64    `db:"order_number" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}
