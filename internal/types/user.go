package types

// User represents a user in the system.
// It is used for handling user-related requests and responses.
type User struct {
	UserId int    `json:"user_id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Age    int    `json:"age"`
}
