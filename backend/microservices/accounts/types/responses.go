package types

import "accounts/dal"

type AccountsResponse struct {
	Accounts []dal.Account `json:"accounts"`
}
