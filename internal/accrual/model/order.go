package model

type Order struct {
	Number  uint64  `db:"number" json:"number" required:"true"`
	Status  string  `db:"status" json:"status"`
	Accrual *uint64 `db:"accrual" json:"accrual"`
}
