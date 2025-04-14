package dal

import (
	"time"
)

type AccountLimits struct {
	Withdrawal int `json:"withdrawal"`
}

type Account struct {
	AccountNumber int
	AccountType string

	Currency string
	Balance int

	Limits AccountLimits
	Status string
	CreatedAt time.Time
}
