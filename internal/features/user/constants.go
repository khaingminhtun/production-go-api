package user

// Role represents the role of a user in the system.
type Role string

// Status represents the lifecycle status of a user account.
type Status string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

const (
	StatusActive   Status = "active"
	StatusInactive Status = "inactive"
	StatusBanned   Status = "banned"
)
