package controllers

import (
	"auth"
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

		if ok := auth.VerifyAuthRules(errors, auth.AuthRules{
			Auth:              true,
			EmailConfirmation: true,
			KYCVerified:       false,
		}, *authInfo); !ok {
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
