package controllers

import (
	"encoding/json"
	"middlewares"
	"net/http"
	"os"
	"strconv"

	"github.com/golang-jwt/jwt"

	"user/dal"
	"user/error"
	"user/types"
	"user/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		services := middlewares.GetContext(r.Context())
		w.Header().Add("Content-Type", "application/json")
		errors := error.NewError(true, w)
		var body types.SignUpRequest
		// TODO: Implement error validation
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errors.ThrowError()
			return
		}
		defer r.Body.Close()

		passwordLeaked, passwordLeakedTimes := utils.PasswordLeaked(body.Password)
		if passwordLeaked {
			errors.NewError("Your password has been leaked: " + strconv.Itoa(passwordLeakedTimes) + " times.")
		}
		if !utils.VerifyPasswordStrength(body.Password) {
			errors.NewError("Password must be at least 8 characters long")
		}

		if !utils.IsValidEmail(body.Email) {
			errors.NewError("Invalid email address")
		}

		emailInUse, err := dal.EmailInUse(services.DB, body.Email)

		if emailInUse && err == nil {
			errors.NewError("Email already in use")
		} else if err != nil {
			errors.ThrowInternalError()
			return
		}

		if len(body.Username) < 5 {
			errors.NewError("Username must be at least 5 characters long")
		}

		usernameInUse, err := dal.UsernameInUse(services.DB, body.Username)

		if usernameInUse && err == nil {
			errors.NewError("Username already in use")
		} else if err != nil {
			errors.ThrowInternalError()
			return
		}

		if errors.ErrorsExist() {
			errors.ThrowError()
			return
		} else {
			var response types.SignUpResponse
			response.Successful = true

			passwordHash, err := utils.HashPassword(body.Password)
			if err != nil {
				errors.ThrowInternalError()
				return
			}

			id, err := dal.CreateUser(services.DB, body.Username, passwordHash, body.Email)
			if err != nil {
				errors.ThrowInternalError()
				return
			}

			claims := jwt.MapClaims{
				"userId": id,
			}

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
			key := []byte(os.Getenv("JWT_SECRET"))
			signedToken, err := token.SignedString(key)

			if err != nil {
				errors.ThrowInternalError()
				return
			}

			w.Header().Add("Authorization", "Bearer "+signedToken)
			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(404)
		return
	}
}
