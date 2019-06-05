package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/kennedymj97/test/todo"
	"golang.org/x/crypto/bcrypt"
)

type UserHandler struct {
	*httprouter.Router
	UserService todo.UserService
	Logger      *log.Logger
}

func NewUserHandler() *UserHandler {
	h := &UserHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.POST("/api/users/create", h.handleCreateUser)
	h.POST("/api/users/login", h.handleLogin)
	h.DELETE("/api/users/logout", h.handleLogout)
	h.DELETE("/api/users/delete", h.handleDeleteUser)
	return h
}

type createUserRequest struct {
	Email    todo.Email `json:"email"`
	Password string     `json:"password"`
}

func (h *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		Error(w, todo.Error("failed to hash password"), http.StatusInternalServerError, h.Logger)
	}

	switch err := h.UserService.CreateUser(req.Email, string(hash)); err {
	case nil:
		encodeJSON(w, &infoResponse{fmt.Sprintf("User has been created with email: %s", req.Email)}, h.Logger)
	case todo.ErrEmailRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case todo.ErrPasswordRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case todo.ErrEmailExists:
		Error(w, err, http.StatusConflict, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type loginRequest struct {
	Email    todo.Email `json:"email"`
	Password string     `json:"password"`
}

func (h *UserHandler) handleLogin(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	userID, pword, err := h.UserService.User(req.Email)
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(pword), []byte(req.Password))
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}

	//generate session id
	sessionID := todo.SessionID(uuid.New().String())
	//get expiry time
	expiryTime := todo.ExpiryTime(time.Now().String())
	switch err := h.UserService.CreateUserSession(sessionID, userID, expiryTime); err {
	case nil:
		// Set to secure in future so it can only be transferred over https
		sessionCookie := &http.Cookie{
			Name:     "session",
			Value:    string(sessionID),
			HttpOnly: true,
			Expires:  time.Now().Add(24 * 14 * time.Hour),
			Path:     "/",
			Domain:   r.Host,
		}
		http.SetCookie(w, sessionCookie)
		encodeJSON(w, &infoResponse{"Login successful"}, h.Logger)
	case todo.ErrSessionRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case todo.ErrUserIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

func (h *UserHandler) handleLogout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionIDCookie, err := r.Cookie("session")
	if err != nil {
		Error(w, err, http.StatusBadRequest, h.Logger)
		return
	}
	sessionID := sessionIDCookie.Value
	switch err := h.UserService.LogoutUser(todo.SessionID(sessionID)); err {
	case nil:
		// Set to secure in future so it can only be transferred over https
		sessionCookie := &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			Expires:  time.Now(),
			Path:     "/api/",
			Domain:   r.Host,
		}
		http.SetCookie(w, sessionCookie)
		encodeJSON(w, &infoResponse{"Succesfully logged out"}, h.Logger)
	case todo.ErrSessionRequired:
		Error(w, todo.ErrSessionRequired, http.StatusBadRequest, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

func (h *UserHandler) handleDeleteUser(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch err := h.UserService.DeleteUser(todo.UserID(r.Header.Get("userID"))); err {
	case nil:
		// Set to secure in future so it can only be transferred over https
		sessionCookie := &http.Cookie{
			Name:     "session",
			Value:    "",
			HttpOnly: true,
			Expires:  time.Now(),
			Path:     "/api/",
			Domain:   r.Host,
		}
		http.SetCookie(w, sessionCookie)
		encodeJSON(w, &infoResponse{"User deleted successfully"}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}
