package postgres

import (
	"github.com/kennedymj97/todo-api"
)

var _ todo.TaskService = &TaskService{}

type TaskService struct {
	client *Client
}

func (s *TaskService) Tasks(id todo.UserID) (*todo.Tasks, error) {
	tx, err := s.client.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()
	rows, err := tx.Query("SELECT taskID, content, completed, timestamp FROM todo.tasks WHERE userID=$1", id)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()
	var todos todo.Tasks
	for rows.Next() {
		tempTask := &todo.Task{}
		if err := rows.Scan(&tempTask.ID, &tempTask.Content, &tempTask.Completed, &tempTask.Timestamp); err != nil {
			tx.Rollback()
			return nil, err
		}
		todos = append(todos, *tempTask)
	}
	return &todos, nil
}

func (s *TaskService) CreateTask(taskID todo.TaskID, content todo.TaskContent, userID todo.UserID) error {
	if FormatInput(content) == "" {
		return todo.ErrTaskContentRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("INSERT INTO todo.tasks(taskID, userID, content) VALUES($1, $2, $3)", taskID, userID, content)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *TaskService) EditTaskStatus(id todo.TaskID, val bool) error {
	if FormatInput(id) == "" {
		return todo.ErrTaskIDRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("UPDATE todo.tasks SET completed=$1 WHERE taskID=$2", val, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *TaskService) ToggleAll(val bool) error {
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("UPDATE todo.tasks SET completed=$1", val)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *TaskService) EditTask(id todo.TaskID, newContent todo.TaskContent) error {
	if FormatInput(id) == "" {
		return todo.ErrTaskIDRequired
	} else if FormatInput(newContent) == "" {
		return todo.ErrTaskContentRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("UPDATE todo.tasks SET content=$1 WHERE taskID=$2", newContent, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *TaskService) DeleteTask(id todo.TaskID) error {
	if FormatInput(id) == "" {
		return todo.ErrTaskIDRequired
	}
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("DELETE FROM todo.tasks WHERE taskid=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (s *TaskService) ClearCompleted() error {
	tx, err := s.client.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Commit()
	_, err = tx.Exec("DELETE FROM todo.tasks WHERE completed=true")
	if err != nil {
		return err
	}
	return nil
}
