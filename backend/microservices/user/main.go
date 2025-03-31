package main

import (
	"db"
	"middlewares"
	"net/http"

	"user/controllers"
)

var con = db.Connect("admin", "admin", "bank")
var ctxMiddleware = middlewares.GetMiddleware(con)

func main() {
	http.Handle("/", ctxMiddleware(http.HandlerFunc(controllers.HomeController)))

	http.ListenAndServe(":8080", nil)
}
