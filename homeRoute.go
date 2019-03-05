package main

import "net/http"

// HomeHandler serves the home page for users to register
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./static/index.html")
}
