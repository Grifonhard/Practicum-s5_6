package model

import "time"

type Good struct {
	ID          uint64    `db:"id" json:"-"`
	Description string    `db:"description" json:"description"`
	Price       uint64    `db:"price" json:"price"`
	OrderID     uint64    `db:"order_id" json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
}
