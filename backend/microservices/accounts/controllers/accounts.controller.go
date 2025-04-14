package controllers

import (
	"accounts/dal"
	"accounts/types"
	"encoding/json"
	"error"
	"log"
	"middlewares"
	"net/http"

	"auth"
)

func Accounts(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authInfo := middlewares.GetAuth(r.Context())
		errors := error.NewError(true, w)
		services := middlewares.GetContext(r.Context())

		if ok := auth.VerifyAuthRules(errors, auth.AuthRules{
			Auth:              true,
			EmailConfirmation: true,
			KYCVerified:       true,
		}, *authInfo); !ok {
			return
		}

		accounts, err := dal.GetAccounts(services.DB, authInfo.UserId)
		if err != nil {
			log.Fatal(err)
			errors.ThrowInternalError()
			return
		}

		var res types.AccountsResponse
		res.Accounts = accounts
		if err := json.NewEncoder(w).Encode(res); err != nil {
			log.Fatal(err)
		}
	}
}
