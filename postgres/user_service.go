package postgres

import (
	"github.com/kennedymj97/todo-api"
	"github.com/lib/pq"
)

var _ todo.UserService = &UserService{}

type UserService struct {
	client *Client
}

func (s *UserService) CreateUser(email todo.Email, password string) error {
	if FormatInput(email) == "" {
		return todo.ErrEmailRequired
	} else if FormatInput(password) == "" {
		return todo.ErrPasswordRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("INSERT INTO todo.users(email, password) VALUES($1, $2)", email, password)
	if err != nil {
		tx.Rollback()
		switch err.(*pq.Error).Message {
		case `duplicate key value violates unique constraint "users_email_key"`:
			return todo.ErrEmailExists
		}
	}
	return nil
}

func (s *UserService) User(email todo.Email) (todo.UserID, string, error) {
	if FormatInput(email) == "" {
		return "", "", todo.ErrEmailRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return "", "", err
	}
	defer tx.Commit()
	row := tx.QueryRow("SELECT userID, password FROM todo.users WHERE email=$1", email)
	var userId todo.UserID
	var pword string
	if err := row.Scan(&userId, &pword); err != nil {
		tx.Rollback()
		return "", "", err
	}
	return userId, pword, nil
}

func (s *UserService) CreateUserSession(sessionID todo.SessionID, userID todo.UserID, expiryTime todo.ExpiryTime) error {
	if FormatInput(sessionID) == "" {
		return todo.ErrSessionRequired
	} else if FormatInput(userID) == "" {
		return todo.ErrUserIDRequired
	} else if FormatInput(expiryTime) == "" {
		return todo.ErrExpiryTimeRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("INSERT INTO todo.userSessions(sessionID, userID, expiryTime) VALUES($1, $2, $3)", sessionID, userID, expiryTime)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) AuthenticateUser(id todo.SessionID) (todo.UserID, error) {
	if FormatInput(id) == "" {
		return "", todo.ErrSessionRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Commit()
	row := tx.QueryRow("SELECT userID FROM todo.userSessions WHERE sessionID=$1", id)
	var userID todo.UserID
	if err := row.Scan(&userID); err != nil {
		return "", err
	}
	return userID, nil
}

func (s *UserService) LogoutUser(id todo.SessionID) error {
	if FormatInput(id) == "" {
		return todo.ErrSessionRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("DELETE FROM todo.userSessions WHERE sessionID=$1", string(id))
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(id todo.UserID) error {
	if FormatInput(id) == "" {
		return todo.ErrUserIDRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("DELETE FROM todo.users WHERE userID=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM todo.userSessions WHERE userID=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("DELETE FROM todo.tasks WHERE userID=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// func (s *UserSessionService) RefreshExpiryTime(id todo.SessionID, newID todo.SessionID, newExpiryTime todo.ExpiryTime) error {
// 	tx, err := s.client.db.Begin()
// 	if err != nil {
// 		return err
// 	}
// 	_, err = tx.Exec("UPDATE todo.userSessions SET sessionID=$1, expiryTime=$2 WHERE sessionID=$3", newID, newExpiryTime, id)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }
