package controllers

import (
	"crypto/sha512"
	"encoding/hex"
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
		db := middlewares.GetContext(r.Context())
		w.Header().Add("Content-Type", "application/json")
		errors := error.NewError(true, w)
		var body types.SignUpRequest
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errors.ThrowError()
			return
		}

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

		successful, emailInUse := dal.EmailInUse(db, body.Email)

		if emailInUse && successful {
			errors.NewError("Email already in use")
		} else if !successful {
			errors.ThrowInternalError()
			return
		}

		if len(body.Username) < 5 {
			errors.NewError("Username must be at least 5 characters long")
		}

		successful, usernameInUse := dal.UsernameInUse(db, body.Username)

		if usernameInUse && successful {
			errors.NewError("Username already in use")
		} else if !successful {
			errors.ThrowInternalError()
			return
		}

		if errors.ErrorsExist() {
			errors.ThrowError()
			return
		} else {
			var response types.SignUpResponse
			response.Successful = true

			hasher := sha512.New()
			if _, err := hasher.Write([]byte(body.Password)); err != nil {
				errors.ThrowInternalError()
				return
			}

			passwordHash := hex.EncodeToString(hasher.Sum(nil))

			successful, id := dal.CreateUser(db, body.Username, passwordHash, body.Email)
			if !successful {
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
