package main

import (
	"fmt"
	"net/http"
)

// Send responses to the user

func HandlerRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world!")
}

func HandlerHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to GoLang RESTful API!")
}
