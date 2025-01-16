package server

import (
	"net/http"
)

type WebServer struct {
	Router        *http.ServeMux
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {

	return &WebServer{
		Router:        http.NewServeMux(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}

func (s *WebServer) Start() {
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}

	// Wrap the mux with LoggingMiddleware
	loggedRouter := LoggingMiddleware(s.Router)

	http.ListenAndServe(s.WebServerPort, loggedRouter)
}
