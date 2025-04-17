package controllers

import (
	"context"
	"encoding/json"
	"error"
	"fmt"
	"log"
	"middlewares"
	"net/http"
	"os"
	"strconv"
	"time"
	"user/dal"
	"user/mql"
	"user/types"
	"user/utils"

	"github.com/golang-jwt/jwt"
	"github.com/redis/go-redis/v9"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		w.Header().Add("Content-Type", "application/json")
		services := middlewares.GetContext(r.Context())
		redis := r.Context().Value("redis").(*redis.Client)
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

		userExist, err := dal.UserExist(services.DB, body.Email, passwordHash)
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

		userId, err := dal.GetUserId(services.DB, body.Email)
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

		cmd := redis.Set(context.Background(), strconv.Itoa(userId), signedToken, time.Minute * 10)
		if cmd.Err() != nil {
			log.Fatal(cmd.Err())
		}

		mql.Send2FSignIn(services.Rabbitmq, userId, body.Email)

		//		var res types.SignInResponse
		//		res.Successful = true

		//		w.Header().Add("Authorization", "Bearer "+signedToken)
		//		json.NewEncoder(w).Encode(res)
		fmt.Fprint(w, "2fa req")
	} else {
		w.WriteHeader(404)
	}
}
