package user

import "time"

// UpdateEmailRequest is used by an authenticated user
// to update their account email.
type UpdateEmailRequest struct {
	Email string `json:"email" binding:"required,email,max=255"`
}

// AdminUpdateUserRequest is used by an admin
// to manage another user's account.
type AdminUpdateUserRequest struct {
	Email  string `json:"email" binding:"required,email,max=255"`
	Role   string `json:"role" binding:"omitempty,oneof=user admin"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive suspended"`
}

// UserResponse is the normal public representation
// of a user account.
type UserResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// AdminUserResponse contains account information
// visible to administrators.
type AdminUserResponse struct {
	ID            uint       `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	Status        string     `json:"status"`
	EmailVerified bool       `json:"email_verified"`
	LastLoginAt   *time.Time `json:"last_login_at"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

// UserListResponse is the paginated user list
// returned to administrators.
type UserListResponse struct {
	Users  []AdminUserResponse `json:"users"`
	Total  int64               `json:"total"`
	Offset int                 `json:"offset"`
	Limit  int                 `json:"limit"`
}

// Model -> DTO
func toUserResponse(user *User) *UserResponse {
	if user == nil {
		return nil
	}

	return &UserResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}
}

func toAdminUserResponse(user *User) *AdminUserResponse {
	if user == nil {
		return nil
	}

	return &AdminUserResponse{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		Role:          user.Role,
		Status:        user.Status,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
	}
}
