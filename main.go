package main

// Init the process handler's registration in router
// Handlers are in handlers.go
// Paths registration go from main -> server -> router
func main() {
	server := NewServer(":3000")
	server.Handle("GET", "/", HandlerRoot)
	server.Handle("GET", "/api", server.AddMiddleware(HandlerHome, CheckAuth(), Loggin()))
	server.Handle("POST", "/api", server.AddMiddleware(HandlerHome, CheckAuth(), Loggin()))
	server.Handle("POST", "/user", UserPostRequest)
	server.Listen()
}
