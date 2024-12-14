package server

import "net/http"

func InitServer(port string, handler http.Handler) error {
	server := http.Server{
		Addr:    ":" + port,
		Handler: handler,
	}

	return server.ListenAndServe()
}
