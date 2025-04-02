package controllers

import (
	"fmt"
	"middlewares"
	"net/http"
	"user/dal"
	"user/error"
)

func EmailConfirmation(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		token := r.URL.Query().Get("token")
		error := error.NewError(true, w)
		services := middlewares.GetContext(r.Context())

		if token == "" {
			error.NewError("Invlaid or Expired token")
            error.ThrowError()
            return
		}

		successful, err := dal.VerifyToken(services.DB, token)
		if err != nil {
			error.ThrowInternalError()
			return
		}

		if !successful {
            error.NewError("Invalid or Expired token")
            error.ThrowError()
            return
		}

		if error.ErrorsExist() {
			error.ThrowError()
			return
		} else {
			fmt.Fprint(w, "Nice")
		}

	} else {
		w.WriteHeader(404)
	}
}
