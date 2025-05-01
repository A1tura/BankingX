package controllers

import (
	"accounts/dal"
	"accounts/types"
	"auth"
	"database/sql"
	"encoding/json"
	"error"
	"middlewares"
	"net/http"
)

func UnfrozeAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		authInfo := middlewares.GetAuth(r.Context())
		services := middlewares.GetContext(r.Context())

		errors := error.NewError(true, w)

		if ok := auth.VerifyAuthRules(errors, auth.AuthRules{
			Auth:              true,
			EmailConfirmation: true,
			KYCVerified:       true,
		}, *authInfo); !ok {
			return
		}

		var request types.UnfrozeAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			errors.NewError("Invalid request")
			errors.ThrowError()
			return
		}

		accountInfo, err := dal.GetAccount(services.DB, request.AccountNumber)
		if err != nil {
			if err == sql.ErrNoRows {
				errors.NewError("Invalid account number.")
				errors.ThrowError()
				return
			}
		}

		if accountInfo.UserId != authInfo.UserId {
			errors.NewError("Invalid account number.")
			errors.ThrowError()
			return
		}

		if accountInfo.Status == "closed" {
			errors.NewError("Account is closed.")
			errors.ThrowError()
			return
		}

		if accountInfo.Status != "frozen" {
			errors.NewError("Account is not frozen")
			errors.ThrowError()
			return
		}

		err = dal.UnfrozeAccount(services.DB, accountInfo.AccountNumber)
		if err != nil {
			errors.ThrowInternalError()
		}

		var res types.UnfrozeAccountResponse
		res.Successfully = true

		if err := json.NewEncoder(w).Encode(res); err != nil {
			errors.ThrowInternalError()
		}
		return
	} else {
		w.WriteHeader(404)
		return
	}
}
