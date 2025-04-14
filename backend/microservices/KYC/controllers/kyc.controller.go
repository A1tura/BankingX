package controllers

import (
	"encoding/json"
	"error"
	"kyc/dal"
	"kyc/types"
	"middlewares"
	"net/http"

	"auth"
)

func KYC(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		errors := error.NewError(true, w)
		services := middlewares.GetContext(r.Context())
		authInfo := middlewares.GetAuth(r.Context())

		if ok := auth.VerifyAuthRules(errors, auth.AuthRules{
			Auth:              true,
			EmailConfirmation: true,
			KYCVerified:       false,
		}, *authInfo); !ok {
			return
		}

		kycStatus, err := dal.AlreadyVerificated(services.DB, authInfo.UserId)
		if err != nil {
			errors.ThrowInternalError()
			return
		}
		if kycStatus {
			errors.NewError("Your KYC process has already been completed and verified. You cannot submit again.")
			errors.ThrowError()
			return
		}

		documentsUploaded, err := dal.DocumentsUploaded(services.DB, authInfo.UserId)
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		if !documentsUploaded {
			errors.NewError("You cannot submit your KYC at this time. Please ensure all required documents are uploaded before proceeding.")
			errors.ThrowError()
			return
		}

		var request types.KYCRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			errors.NewError("Invalid request")
			errors.ThrowError()
			return
		}

		if err := dal.CreateKYC(services.DB, authInfo.UserId, request.FirstName, request.MiddleName, request.LastName, request.DateOfBirth, request.PhoneNumber, request.IdNumber, request.Country, request.State, request.City, request.Address, request.PostalCode); err != nil {
			errors.ThrowInternalError()
			return
		}

		if errors.ErrorsExist() {
			errors.ThrowError()
			return
		} else {
			var response types.KYCResponse
			response.Successfully = true

			if err := json.NewEncoder(w).Encode(response); err != nil {
				errors.ThrowInternalError()
				return
			}
		}
	} else {
		w.WriteHeader(404)
		return
	}
}
