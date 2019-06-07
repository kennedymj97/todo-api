package http

import (
	"net/http"
	"time"
)

func InitServer() *http.Server {
	s := &http.Server{
		Addr:         ":https",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	return s
}
