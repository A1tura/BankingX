package controllers

import (
	"fmt"
	"kyc/dal"
	"middlewares"
	"net/http"
	"strconv"
)

func Approve(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		services := middlewares.GetContext(r.Context())
		params := r.URL.Query()
		userIdStr := params.Get("userId")
		if userIdStr == "" {
			fmt.Fprint(w, "Missing parameter userId")
			w.WriteHeader(400)
			return
		}

		userId, err := strconv.Atoi(userIdStr)
		if err != nil {
			fmt.Fprint(w, "Invalid user id")
			w.WriteHeader(400)
			return
		}

		if err := dal.ApproveKYC(services.DB, userId); err != nil {
			fmt.Fprint(w, err)
			return
		}

		fmt.Fprint(w, "OK")
		return
	}
}
