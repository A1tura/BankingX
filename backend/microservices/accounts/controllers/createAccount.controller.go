package controllers

import (
	"accounts/dal"
	"accounts/types"
	"accounts/utils"
	"auth"
	"encoding/json"
	"error"
	"fmt"
	"middlewares"
	"net/http"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Add("Content-Type", "application/json")
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

		var request types.CreateAccountRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			errors.NewError("Invalid request")
			errors.ThrowError()
			return
		}
		defer r.Body.Close()

		if request.AccountType != "saving" && request.AccountType != "checking" && request.AccountType != "business" {
			errors.NewError("Invalid account type.")
			errors.ThrowError()
			return
		}

		if request.Currency != "USD" && request.Currency != "EUR" {
			errors.NewError("Invalid currency")
			errors.ThrowError()
			return
		}

		var verifiedBankNumber string

		for verifiedBankNumber == "" {
			bankNumber, err := utils.GenerateAccountNumber()
			bankNumber = "0030" + bankNumber
			if err != nil {
				errors.ThrowInternalError()
				return
			}

			bankNumberExist, err := dal.AccountNumberExist(services.DB, bankNumber)

			if !bankNumberExist {
				verifiedBankNumber = bankNumber
			}
		}
		var limits dal.AccountLimits

		_, err := dal.CreateAccount(services.DB, authInfo.UserId, request.AccountType, verifiedBankNumber, request.Currency, request.IsPrimary, limits)
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		fmt.Fprint(w, verifiedBankNumber)
	} else {
		w.WriteHeader(404)
		return
	}
}
