package postgres

import (
	"database/sql"
	"fmt"
	"os"
	"reflect"
	"strings"

	"github.com/joho/godotenv"
	"github.com/kennedymj97/todo-api"
	_ "github.com/lib/pq"
)

type Client struct {
	db          *sql.DB
	taskService TaskService
	userService UserService
}

func NewClient() *Client {
	c := &Client{}
	c.taskService.client = c
	c.userService.client = c
	return c
}

func (c *Client) Open() error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	host, ok := os.LookupEnv("DBHOST")
	if !ok {
		return todo.ErrDBHOSTRequried
	}
	port, ok := os.LookupEnv("DBPORT")
	if !ok {
		return todo.ErrDBPORTRequried
	}
	user, ok := os.LookupEnv("DBUSER")
	if !ok {
		return todo.ErrDBUSERRequried
	}
	pword, ok := os.LookupEnv("DBPASSWORD")
	if !ok {
		return todo.ErrDBPASSWORDRequried
	}
	dbName, ok := os.LookupEnv("DBNAME")
	if !ok {
		return todo.ErrDBNAMERequried
	}
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, pword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}

	err = db.Ping()
	if err != nil {
		return err
	}

	// Code to set up new schema and table if not already set up.
	db.Exec("CREATE SCHEMA IF NOT EXISTS todo;")
	db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
	newTaskTable := `CREATE TABLE IF NOT EXISTS todo.tasks( 
	taskID UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
	userID UUID NOT NULL,
	content TEXT NOT NULL, 
	completed BOOL NOT NULL DEFAULT false, 
	timestamp TIMESTAMP NOT NULL DEFAULT current_timestamp
	);`
	newUserTable := `CREATE TABLE IF NOT EXISTS todo.users(
	userID UUID PRIMARY KEY DEFAULT uuid_generate_v1(),
	email TEXT NOT NULL UNIQUE,
	password TEXT NOT NULL
	);`
	// CHANGE EXPIRY TIME TO BE A CHAR OF FIXED LENGTH
	newUserSessionTable := `CREATE TABLE IF NOT EXISTS todo.userSessions(
	sessionID UUID PRIMARY KEY,
	userID UUID NOT NULL,
	expiryTime TEXT NOT NULL 
	);`
	db.Exec(newTaskTable)
	db.Exec(newUserTable)
	db.Exec(newUserSessionTable)

	c.db = db

	fmt.Println("Succesfully connected to database")
	return nil
}

func (c *Client) Close() error {
	if c.db != nil {
		return c.db.Close()
	}
	return nil
}

func (c *Client) TaskService() todo.TaskService { return &c.taskService }

func (c *Client) UserService() todo.UserService { return &c.userService }

func FormatInput(input interface{}) string {
	s := reflect.ValueOf(input).String()
	return strings.TrimSpace(s)
}
