package main

import (
	"crypto/tls"
	"log"
	goHttp "net/http"
	"time"

	"github.com/kennedymj97/todo-api/http"
	"github.com/kennedymj97/todo-api/postgres"
	"golang.org/x/crypto/acme/autocert"
)

func main() {
	// Connect to db and create services
	dbClient := postgres.NewClient()
	err := dbClient.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer dbClient.Close()

	// Create new handler
	taskHandler := http.NewTaskHandler()
	userHandler := http.NewUserHandler()
	taskHandler.TaskService = dbClient.TaskService()
	userHandler.UserService = dbClient.UserService()

	// Create new server
	certManager := autocert.Manager{
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist("api.mattkennedy.io"),
		Cache:      autocert.DirCache("certs"),
	}

	s := &goHttp.Server{
		Addr:         ":https",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig: &tls.Config{
			GetCertificate: certManager.GetCertificate,
		},
	}
	s.Handler = &http.Handler{TaskHandler: taskHandler, UserHandler: userHandler}

	go goHttp.ListenAndServe(":http", certManager.HTTPHandler(nil))

	log.Fatal(s.ListenAndServeTLS("", ""))
}
