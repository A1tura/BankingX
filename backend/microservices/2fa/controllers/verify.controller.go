package controllers

import (
	"2fa/types"
	"auth"
	"context"
	"encoding/json"
	"error"
	"fmt"
	"middlewares"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func VerifyAuth(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		authInfo := middlewares.GetAuth(r.Context())
		redisCon := r.Context().Value("redis").(*redis.Client)
		authRedisCon := r.Context().Value("redisAuth").(*redis.Client)
		errors := error.NewError(true, w)

		if ok := auth.VerifyAuthRules(errors, auth.AuthRules{
			Auth:              false,
			EmailConfirmation: false,
			KYCVerified:       false,
		}, *authInfo); !ok {
			return
		}

		var req types.VerifyRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			errors.NewError("Invalid request")
			errors.ThrowError()
			return
		}

		cmd := redisCon.Get(context.Background(), fmt.Sprintf("auth:%s", req.Code))
		userId, err := cmd.Result()
		if err != nil {
			errors.NewError("Invalid code")
			errors.ThrowError()
			return
		}

		cmd = authRedisCon.Get(context.Background(), userId)
		jwt, err := cmd.Result()
		if err != nil {
			errors.NewError("Invalid code")
			errors.ThrowError()
			return
		}

		authRedisCon.Del(context.Background(), userId)
		redisCon.Del(context.Background(), fmt.Sprintf("auth:%s", req.Code))

		w.Header().Add("Authorization", "Bearer "+jwt)

		var res types.AuthVerificationResponse
		res.Successfully = true

		json.NewEncoder(w).Encode(res)
	} else {
		w.WriteHeader(404)
	}
}
