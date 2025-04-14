package auth

import (
	"error"
	"middlewares"
)

type AuthRules struct {
	Auth              bool
	EmailConfirmation bool
	KYCVerified       bool
}

func VerifyAuthRules(errors error.Error, authRules AuthRules, authInfo middlewares.AuthInfo) bool {
	if authRules.Auth && !authInfo.IsAuth {
		errors.NewError("You must be authenticated to access this resource.")
		errors.ThrowError()
		return false
	}

	if authRules.EmailConfirmation && authInfo.EmailConfirmed == nil {
		errors.ThrowInternalError()
		return false
	}

	if authRules.EmailConfirmation && !*authInfo.EmailConfirmed {
		errors.NewError("Your email address is not yet confirmed. Please verify your email before accessing this resource.")
		errors.ThrowError()
		return false
	}

	if authRules.KYCVerified && authInfo.KYCStatus == nil {
		errors.ThrowInternalError()
		return false
	}

	if authRules.KYCVerified && *authInfo.KYCStatus != "verified" {
		errors.NewError("Access denied. KYC verification required to access this resource.")
		errors.ThrowError()
		return false
	}

	return true
}
