package model

type Good struct {
	Description string `db:"description" json:"description" required:"true"`
	Price       uint64 `db:"price" json:"price" required:"true"`
	OrderID     uint64 `db:"order_id" json:"-"`
}
