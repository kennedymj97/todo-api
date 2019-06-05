package todo

// Error represents a todo error
type Error string

// Error returns the error message
func (e Error) Error() string { return string(e) }

// General errors
const (
	ErrInternal     = Error("internal error")
	ErrUnauthorized = Error("user is not authorized")
)

// Database errors
const (
	ErrDBHOSTRequried     = Error("DBHOST environment variable required but not set")
	ErrDBPORTRequried     = Error("DBPORT environment variable required but not set")
	ErrDBUSERRequried     = Error("DBUSER environment variable required but not set")
	ErrDBPASSWORDRequried = Error("DBPASSWORD environment variable required but not set")
	ErrDBNAMERequried     = Error("DBNAME environment variable required but not set")
)

// http errors
const (
	ErrInvalidJSON = Error("invalid json")
)

// Task errors
const (
	ErrTaskContentRequired   = Error("task content requried")
	ErrTaskIDRequired        = Error("task id required")
	ErrTaskNotFound          = Error("task not found")
	ErrCompletedBoolRequired = Error("completed bool requried")
)

// User errors
const (
	ErrEmailRequired      = Error("email required")
	ErrPasswordRequired   = Error("password required")
	ErrSessionRequired    = Error("session requried")
	ErrExpiryTimeRequired = Error("expiry time required")
	ErrUserIDRequired     = Error("user id requried")
	ErrUsernameExists     = Error("username is taken")
	ErrEmailExists        = Error("email already exists")
)
