package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"

	"github.com/julienschmidt/httprouter"
	"github.com/kennedymj97/todo-api"
)

type TaskHandler struct {
	*httprouter.Router
	TaskService todo.TaskService
	Logger      *log.Logger
}

func NewTaskHandler() *TaskHandler {
	h := &TaskHandler{
		Router: httprouter.New(),
		Logger: log.New(os.Stderr, "", log.LstdFlags),
	}
	h.GET("/api/tasks", h.handleTasks)
	h.POST("/api/tasks/create", h.handleCreateTask)
	h.POST("/api/tasks/edit", h.handleTaskEdit)
	h.POST("/api/tasks/toggle", h.handleTaskToggle)
	h.POST("/api/tasks/toggleAll", h.handleToggleAll)
	h.DELETE("/api/tasks/delete/:id", h.handleDeleteTask)
	h.DELETE("/api/tasks/clearCompleted", h.handleClearCompleted)
	return h
}

type getTasksResponse struct {
	Tasks *todo.Tasks `json:"tasks,omitempty"`
}

func (h *TaskHandler) handleTasks(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	t, err := h.TaskService.Tasks(todo.UserID(r.Header.Get("userID")))
	tasks := *t
	sort.Slice(tasks, func(i, j int) bool {
		return tasks[i].Timestamp < tasks[j].Timestamp
	})
	if err != nil {
		Error(w, err, http.StatusInternalServerError, h.Logger)
	} else if t == nil {
		NotFound(w)
	} else {
		encodeJSON(w, &getTasksResponse{Tasks: &tasks}, h.Logger)
	}
}

type createTaskRequest struct {
	ID      todo.TaskID      `json:"id"`
	Content todo.TaskContent `json:"content,omitempty"`
}

func (h *TaskHandler) handleCreateTask(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req createTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	content := req.Content

	switch err := h.TaskService.CreateTask(req.ID, content, todo.UserID(r.Header.Get("userID"))); err {
	case nil:
		encodeJSON(w, &infoResponse{fmt.Sprintf("Task has been successfully created with content: %s", content)}, h.Logger)
	case todo.ErrTaskContentRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type editTaskRequest struct {
	ID      todo.TaskID      `json:"id"`
	Content todo.TaskContent `json:"content"`
}

func (h *TaskHandler) handleTaskEdit(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req editTaskRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	switch err := h.TaskService.EditTask(req.ID, req.Content); err {
	case nil:
		encodeJSON(w, &infoResponse{fmt.Sprintf("Task has been updated to content: %s", req.Content)}, h.Logger)
	case todo.ErrTaskIDRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	case todo.ErrTaskContentRequired:
		Error(w, err, http.StatusBadRequest, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}

}

type taskStatusRequest struct {
	ID  todo.TaskID `json:"id"`
	Val bool        `json:"val"`
}

func (h *TaskHandler) handleTaskToggle(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Decode request
	var req taskStatusRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}

	// Create task
	switch err := h.TaskService.EditTaskStatus(req.ID, req.Val); err {
	case nil:
		encodeJSON(w, &infoResponse{fmt.Sprintf("Task status has been set to %t", req.Val)}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type toggleAllRequest struct {
	Val bool `json:"val"`
}

func (h *TaskHandler) handleToggleAll(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req toggleAllRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		Error(w, todo.ErrInvalidJSON, http.StatusBadRequest, h.Logger)
		return
	}
	switch err := h.TaskService.ToggleAll(req.Val); err {
	case nil:
		encodeJSON(w, &infoResponse{"Tasks have all been toggled."}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

type deleteTaskRequest struct {
	ID todo.TaskID `json:"id"`
}

func (h *TaskHandler) handleDeleteTask(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	switch err := h.TaskService.DeleteTask(todo.TaskID(p.ByName("id"))); err {
	case nil:
		encodeJSON(w, &infoResponse{"Task has been successfully deleted"}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}

func (h *TaskHandler) handleClearCompleted(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	switch err := h.TaskService.ClearCompleted(); err {
	case nil:
		encodeJSON(w, &infoResponse{"Completed tasks have been succesfully deleted"}, h.Logger)
	default:
		Error(w, err, http.StatusInternalServerError, h.Logger)
	}
}
