package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	_ "middlewares"
	"net/http"
	"user/error"
	"user/types"
	"user/utils"
)

func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
        w.Header().Add("Content-Type", "application/json")
		errors := error.NewError(true, w)
		var body types.SignUpRequest
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(500)
			fmt.Fprint(w, "Error, try again later")
			return
		}

		defer r.Body.Close()
		if err := json.Unmarshal(bodyBytes, &body); err != nil {
			w.WriteHeader(400)
			fmt.Fprint(w, "Bad request")
			return
		}

		if !utils.IsValidEmail(body.Email) {
			errors.NewError("Invalid email address")
		}

		if errors.ErrorsExist() {
			errors.ThrowError()
			return
		} else {
			var response types.SignUpResponse
			response.Successful = true
			w.Header().Add("Authorization", "Bearer xdd")
			fmt.Print(body)

			json.NewEncoder(w).Encode(response)
		}
	} else {
		w.WriteHeader(404)
		return
	}
}
