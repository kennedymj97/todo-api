package http

import (
	"net/http"
	"time"
)

func InitServer() *http.Server {
	s := &http.Server{
		Addr:         ":8080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return s
}
