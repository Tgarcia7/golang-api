package main

import (
	"net/http"
)

// The Router implementation requires ServeHTTP func
type Router struct {
	rules map[string]map[string]http.HandlerFunc // HTTP rules mapping
}

func newRouter() *Router {
	return &Router{
		rules: make(map[string]map[string]http.HandlerFunc),
	}
}

func (router *Router) FindHanlder(path string, method string) (http.HandlerFunc, bool, bool) {
	_, exists := router.rules[path]
	handler, methodExists := router.rules[path][method]
	return handler, methodExists, exists
}

func (router *Router) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	handler, methodExists, exists := router.FindHanlder(request.URL.Path, request.Method)

	// Route not found 404
	if !exists {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !methodExists {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Call the handler (from handlers.go) to attend the request
	handler(w, request)
}
