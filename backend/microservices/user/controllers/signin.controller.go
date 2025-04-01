package controllers

import (
	"encoding/json"
	"middlewares"
	"net/http"
	"os"
	"user/dal"
	"user/error"
	"user/types"
	"user/utils"

	"github.com/golang-jwt/jwt"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Add("Content-Type", "application/json")
		db := middlewares.GetContext(r.Context())
		errors := error.NewError(true, w)

		var body types.SignInRequest
		// TODO: Implement error validation
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			errors.NewError("Invalid form data")
			errors.ThrowError()
			return
		}

		passwordHash, err := utils.HashPassword(body.Password)
		// TODO: Implement error validation
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		userExist, err := dal.UserExist(db, body.Email, passwordHash)
		// TODO: Implement error validation
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		if !userExist {
			errors.NewError("User with that password and email do not exist")
			errors.ThrowError()
			return
		}

		userId, err := dal.GetUserId(db, body.Email)
		// TODO: Implement error validation
		if err != nil {
			errors.ThrowInternalError()
			return
		}

		claims := jwt.MapClaims{
			"userId": userId,
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		key := []byte(os.Getenv("JWT_SECRET"))
		signedToken, err := token.SignedString(key)

		var res types.SignInResponse
		res.Successful = true

		w.Header().Add("Authorization", "Bearer "+signedToken)
		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(404)
	}
}
