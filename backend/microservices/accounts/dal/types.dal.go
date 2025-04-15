package dal

import (
	"time"
)

type AccountLimits struct {
	Withdrawal int `json:"withdrawal"`
}

type Account struct {
	AccountNumber string
	AccountType   string

	Currency string
	Balance  float64

	Limits    AccountLimits
	Status    string
	CreatedAt time.Time
}

type FullAccount struct {
	Id            int
	UserId        int
	AccountNumber string
	AccountType   string

	Currency string
	Balance  float64

	Limits    AccountLimits
	Status    string
	UpdatedAt time.Time
	CreatedAt time.Time
}
