package model

import "time"

type AccrualProgram struct {
	ID         uint64    `db:"id" json:"-"`
	Match      string    `db:"match" json:"match"`
	Reward     int64     `db:"reward" json:"reward"`
	RewardType string    `db:"reward_type" json:"reward_type"`
	CreatedAt  time.Time `db:"created_at" json:"-"`
}
