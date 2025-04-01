package controllers

import (
	"fmt"
	"net/http"
)

func SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
        fmt.Fprint(w, "/signIn")
	} else {
        w.WriteHeader(404)
    }
}
