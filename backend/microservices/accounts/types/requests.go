package types

type CreateAccountRequest struct {
	AccountType string `json:"accountType"`
	Currency    string `json:"currency"`
	IsPrimary   bool   `json:"isPrimary"`
}

type FrozeAccountRequest struct {
	AccountNumber string `json:"accountNumber"`
}
