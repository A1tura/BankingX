package controllers

import (
	"accounts/dal"
	"accounts/types"
	"auth"
	"database/sql"
	"encoding/json"
	"error"
	"fmt"
	"log"
	"middlewares"
	"net/http"
)

func FrozeAccount(w http.ResponseWriter, r *http.Request) {
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

		var request types.FrozeAccountRequest

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

		if accountInfo.Status == "frozen" {
			errors.NewError("Account is already frozen.")
			errors.ThrowError()
			return
		}

		if err := dal.FrozeAccount(services.DB, accountInfo.AccountNumber); err != nil {
			log.Fatal(err)
		}

		fmt.Fprint(w, accountInfo)

	} else {
		w.WriteHeader(404)
		return
	}
}
