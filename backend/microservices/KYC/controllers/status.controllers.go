package controllers

import (
	"encoding/json"
	"error"
	"kyc/dal"
	"kyc/types"
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

		if authInfo.EmailConfirmed == nil {
			errors.ThrowInternalError()
			return
		}

		if !*authInfo.EmailConfirmed {
			errors.NewError("Your email address is not yet confirmed. Please verify your email before accessing this resource.")
			errors.ThrowError()
			return
		}

		KYCStatus, err := dal.KYCStatus(services.DB, authInfo.UserId)
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		if KYCStatus == "NE" {
			errors.NewError("KYC status unavailable: no verification request has been submitted.")
			errors.ThrowError()
			return
		}

		var response types.KYCStatutsResponse

		response.Successfully = true
		response.Status = KYCStatus

		if err := json.NewEncoder(w).Encode(response); err != nil {
			errors.ThrowInternalError()
			return
		}

	} else {
		w.WriteHeader(404)
		return
	}
}
