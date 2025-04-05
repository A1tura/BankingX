package controllers

import (
	"error"
	"fmt"
	"kyc/dal"
	"middlewares"
	"net/http"
)

func Status(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		errors := error.NewError(true, w)
		authInfo := middlewares.GetAuth(r.Context())
		services := middlewares.GetContext(r.Context())

		if !authInfo.IsAuth {
			errors.NewError("You must be authenticated to access this resource.")
			errors.ThrowError()
			return
		}

		emailConfirmed, err := dal.EmailConfirmed(services.DB, authInfo.UserId)
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		if !emailConfirmed {
			errors.NewError("Your email address is not yet confirmed. Please verify your email before accessing this resource.")
			errors.ThrowError()
			return
		}

        KYCStatus, err := dal.KYCStatus(services.DB, authInfo.UserId)
        if err != nil {
            errors.ThrowInternalError()
            return
        }

        fmt.Fprint(w, KYCStatus)
	} else {
		w.WriteHeader(404)
		return
	}
}
