package controllers

import (
	"fmt"
	"net/http"
)

func EmailConfirmation(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        fmt.Fprint(w, "/emailConfirmation endpoint")
    } else {
        w.WriteHeader(404)
    }
}
