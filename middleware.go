package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func CheckAuth() Middleware {
	return func(nextMiddleware http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, request *http.Request) {

			authenticated := true
			fmt.Println("Checking authentication")

			if authenticated {
				nextMiddleware(w, request)
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		}
	}
}

func Loggin() Middleware {
	return func(nextMiddleware http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {

			start := time.Now()
			defer func() {
				log.Println(r.URL.Path, time.Since(start))
			}()

			nextMiddleware(w, r)

		}
	}
}
