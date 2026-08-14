package auth

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=20"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type RegisterResponse struct {
	RegistrationID string `json:"registration_id" binding:"required"`
	Message        string `json:"message"`
}

type PendingRegistration struct {
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	OTPHash      string    `json:"otp_hash"`
	CreatedAt    time.Time `json:"created_at"`
}

type VerifyRegisterRequest struct {
	RegistrationID string `json:"registration_id" binding:"required"`
	OTP            string `json:"otp" binding:"required"`
}

type VerifyRegisterResponse struct {
	Message string `json:"message"`
}
