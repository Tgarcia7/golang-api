package main

import (
	"net/http"
)

// The Router implementation requires ServeHTTP func
type Router struct {
	rules map[string]http.HandlerFunc // HTTP rules mapping
}

func newRouter() *Router {
	return &Router{
		rules: make(map[string]http.HandlerFunc),
	}
}

func (router *Router) FindHanlder(path string) (http.HandlerFunc, bool) {
	handler, exists := router.rules[path]
	return handler, exists
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler, exists := router.FindHanlder(r.URL.Path)

	// Route not found 404
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	// Call the handler (from handlers.go) to attend the request
	handler(w, r)
}
