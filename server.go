package main

import "net/http"

// Struct properties
type Server struct {
	port   string
	router *Router
}

// Server init
func NewServer(port string) *Server {
	// Exports the server instance, avoid creating more instances
	return &Server{
		port:   port,
		router: newRouter(), // Router instance to handle requests
	}
}

func (server *Server) Handle(method string, path string, handler http.HandlerFunc) {
	_, exists := server.router.rules[path]

	if !exists {
		server.router.rules[path] = make(map[string]http.HandlerFunc)
	}

	server.router.rules[path][method] = handler
}

func (server *Server) Listen() error {
	// Routes main endpoint registration
	// Makes the router start attending routes
	http.Handle("/", server.router)
	// Init server listening
	err := http.ListenAndServe(server.port, nil)

	if err != nil {
		return err
	}

	return nil
}

// Creates the middleware chaining. With ... indicates that we do not know the number of middlewares
func (server *Server) AddMiddleware(middleware http.HandlerFunc, middlewares ...Middleware) http.HandlerFunc {
	// Pass parameters between middlewares
	for _, m := range middlewares {
		middleware = m(middleware)
	}

	return middleware
}
