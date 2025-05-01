package types

import "accounts/dal"

type AccountsResponse struct {
	Accounts []dal.Account `json:"accounts"`
}

type FrozeAccountResponse struct {
	Successfully bool `json:"successfully"`
}


type UnfrozeAccountResponse struct {
	Successfully bool `json:"successfully"`
}
