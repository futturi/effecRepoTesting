package server

import (
	"net/http"
	"time"
)

func NewServer(handler http.Handler, port string) error {
	server := http.Server{
		Addr:         ":" + port,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      handler,
	}
	return server.ListenAndServe()
}
