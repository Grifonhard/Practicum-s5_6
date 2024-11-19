package model

import "time"

type AccrualProgram struct {
	ID         uint64    `db:"id" json:"-"`
	Match      string    `db:"match" json:"match" required:"true"`
	Reward     float64   `db:"reward" json:"reward" required:"true"`
	RewardType string    `db:"reward_type" json:"reward_type" required:"true"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
}
