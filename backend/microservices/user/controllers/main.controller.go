package controllers

import (
	"fmt"
	"middlewares"
	"net/http"
)

func HomeController(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello, World!")
}
