package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/kennedymj97/todo-api"
)

type Handler struct {
	TaskHandler *TaskHandler
	UserHandler *UserHandler
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Implement middleware here
	h.TaskHandler.Logger.Printf("%s %s %s", r.Proto, r.Method, r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8080/")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}
	switch r.URL.Path {
	case "/api/users/login":
		break
	case "/api/users/create":
		break
	default:
		userID, err := h.auth(r)
		if err != nil {
			Error(w, todo.ErrUnauthorized, http.StatusUnauthorized, h.UserHandler.Logger)
			return
		}
		r.Header.Add("userID", string(userID))
	}
	if strings.HasPrefix(r.URL.Path, "/api/tasks") {
		h.TaskHandler.ServeHTTP(w, r)
	} else if strings.HasPrefix(r.URL.Path, "/api/users") {
		h.UserHandler.ServeHTTP(w, r)
	} else {
		http.NotFound(w, r)
	}
}

func (h *Handler) auth(r *http.Request) (todo.UserID, error) {
	sessionCookie, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	sessionID := sessionCookie.Value
	// NEED TO CHECK THAT THE SESSION COOKIE HAS NOT EXPIRED
	// WAY TO REFRESH THE SESSION COOKIE IMPLEMENTED HERE
	userID, err := h.UserHandler.UserService.AuthenticateUser(todo.SessionID(sessionID))
	if err != nil {
		return "", err
	}
	return userID, nil
}

type infoResponse struct {
	Info string `json:"info,omitempty"`
}

type errorResponse struct {
	Err string `json:"err,omitempty"`
}

func Error(w http.ResponseWriter, err error, code int, logger *log.Logger) {
	// Log error
	logger.Printf("http error: %s (code=%d)", err, code)

	// Hide error from client if it is internal
	if code == http.StatusInternalServerError {
		err = todo.ErrInternal
	}

	w.WriteHeader(code)
	json.NewEncoder(w).Encode(&errorResponse{Err: err.Error()})
}

func NotFound(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte(`{}` + "\n"))
}

func encodeJSON(w http.ResponseWriter, v interface{}, logger *log.Logger) {
	if err := json.NewEncoder(w).Encode(v); err != nil {
		Error(w, err, http.StatusInternalServerError, logger)
	}
}
