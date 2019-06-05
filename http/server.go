package http

import (
	"net/http"
	"time"
)

const DefaultAddr = ":8080"

func InitServer() *http.Server {
	s := &http.Server{
		Addr:         DefaultAddr,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}
	return s
}
