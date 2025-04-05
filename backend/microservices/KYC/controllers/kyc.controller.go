package controllers

import (
	"error"
	"fmt"
	"middlewares"
	"net/http"
)

func KYC(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		errors := error.NewError(true, w)
		authInfo := middlewares.GetAuth(r.Context())

        if !authInfo.IsAuth {
            errors.NewError("You must be authenticated to access this resource.")
            errors.ThrowError();
            return
        }

        fmt.Fprint(w, authInfo.UserId)

		if errors.ErrorsExist() {
			errors.ThrowError()
			return
		} else {
			fmt.Fprint(w, "Done")
		}
	} else {
		w.WriteHeader(404)
		return
	}
}
