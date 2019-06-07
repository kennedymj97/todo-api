package main

import (
	"log"

	"github.com/kennedymj97/todo-api/http"
	"github.com/kennedymj97/todo-api/postgres"
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
	s := http.InitServer()
	s.Handler = &http.Handler{TaskHandler: taskHandler, UserHandler: userHandler}
	log.Fatal(s.ListenAndServe())
}
