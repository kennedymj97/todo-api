package todo

type UserID string
type Email string
type SessionID string
type ExpiryTime string

type UserService interface {
	CreateUser(email Email, password string) error
	User(email Email) (UserID, string, error)
	CreateUserSession(sessionId SessionID, userId UserID, expiryTime ExpiryTime) error
	AuthenticateUser(id SessionID) (UserID, error)
	LogoutUser(id SessionID) error
	//RefreshExpiryTime(id SessionID, newId SessionID, newExpiryTime ExpiryTime) error
	DeleteUser(id UserID) error
	//UpdateUser(username Username, email Email, password Password) error
}

type TaskID string
type TaskContent string

type Task struct {
	ID        TaskID      `json:"id"`
	Content   TaskContent `json:"content"`
	Completed bool        `json:"completed"`
	Timestamp string      `json:"timestamp"`
}

type Tasks []Task

type TaskService interface {
	Tasks(id UserID) (*Tasks, error)
	CreateTask(taskID TaskID, content TaskContent, userID UserID) error
	EditTaskStatus(id TaskID, val bool) error
	ToggleAll(val bool) error
	EditTask(id TaskID, newContent TaskContent) error
	DeleteTask(id TaskID) error
	ClearCompleted() error
}
