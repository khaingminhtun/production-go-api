package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/khaingminhtun/production-go-api/internal/features/user"
	redisinfra "github.com/khaingminhtun/production-go-api/internal/infrastructure/redis"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorhandler/apperror"
	"github.com/khaingminhtun/production-go-api/internal/shared/security"
)

type Service interface {
	Register(
		ctx context.Context,
		req RegisterRequest,
	) (*RegisterResponse, error)

	VerifyRegister(
		ctx context.Context,
		req VerifyRegisterRequest,
	) (*VerifyRegisterResponse, error)
}

type service struct {
	userRepo   user.Repository
	redisStore redisinfra.RedisStore
	emailQueue redisinfra.EmailQueue
}

func NewService(
	userRepo user.Repository,
	redisStore redisinfra.RedisStore,
	emailQueue redisinfra.EmailQueue,
) Service {
	return &service{
		userRepo:   userRepo,
		redisStore: redisStore,
		emailQueue: emailQueue,
	}
}

func (s *service) Register(
	ctx context.Context,
	req RegisterRequest,
) (*RegisterResponse, error) {

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	// ============================================================
	// Check existing email
	// ============================================================

	_, err := s.userRepo.GetByEmail(ctx, email)

	switch {
	case err == nil:
		return nil, apperror.New(
			apperror.CodeUserAlreadyExists,
			"user with this email already exists",
			nil,
		)

	case apperror.Is(err, apperror.CodeUserNotFound):
		// Expected.
		// Email does not exist, so registration can continue.

	default:
		return nil, err
	}

	// ============================================================
	// Check existing username
	// ============================================================

	_, err = s.userRepo.GetByUsername(ctx, username)

	switch {
	case err == nil:
		return nil, apperror.New(
			apperror.CodeUserAlreadyExists,
			"user with this username already exists",
			nil,
		)

	case apperror.Is(err, apperror.CodeUserNotFound):
		// Expected.
		// Username does not exist, so registration can continue.

	default:
		return nil, err
	}

	// ============================================================
	// Hash password
	// ============================================================

	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// ============================================================
	// Generate OTP
	// ============================================================

	otp, err := security.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	otpHash := security.HashOTP(otp)

	// ============================================================
	// Build pending registration
	// ============================================================

	pending := PendingRegistration{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		OTPHash:      otpHash,
	}

	data, err := json.Marshal(pending)
	if err != nil {
		return nil, fmt.Errorf(
			"marshal pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Save pending registration in Redis
	// ============================================================
	registrationID := uuid.NewString()

	key := "auth:register:" + registrationID

	const registrationTTL = 10 * time.Minute

	if err := s.redisStore.Set(
		ctx,
		key,
		string(data),
		registrationTTL,
	); err != nil {
		return nil, fmt.Errorf(
			"save pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Create email job
	// ============================================================

	job := redisinfra.EmailJob{
		ID:        uuid.NewString(),
		To:        email,
		Subject:   "Verify your email",
		Template:  "otp_verification",
		CreatedAt: time.Now(),

		Data: map[string]any{
			"OTP":       otp,
			"ExpiresIn": "10 minutes",
		},
	}

	// ============================================================
	// Publish email job
	// ============================================================

	if err := s.emailQueue.Publish(ctx, job); err != nil {

		// Remove pending registration because
		// the email could not be queued.
		_ = s.redisStore.Delete(ctx, key)

		return nil, fmt.Errorf(
			"queue verification email: %w",
			err,
		)
	}

	// ============================================================
	// Response
	// ============================================================

	return &RegisterResponse{
		RegistrationID: registrationID,
		Message:        "Verification code sent to your email",
	}, nil
}

func (s *service) VerifyRegister(
	ctx context.Context,
	req VerifyRegisterRequest,
) (*VerifyRegisterResponse, error) {

	registrationID := strings.TrimSpace(req.RegistrationID)

	key := "auth:register:" + registrationID

	data, err := s.redisStore.Get(ctx, key)
	if err != nil {
		return nil, fmt.Errorf(
			"get pending registration: %w",
			err,
		)
	}

	if data == "" {
		return nil, apperror.New(
			apperror.CodeInvalidRequest,
			"registration expired or not found",
			nil,
		)
	}

	// ============================================================
	// Decode pending registration
	// ============================================================

	var pending PendingRegistration

	if err := json.Unmarshal(
		[]byte(data),
		&pending,
	); err != nil {
		return nil, fmt.Errorf(
			"unmarshal pending registration: %w",
			err,
		)
	}

	// ============================================================
	// Verify OTP
	// ============================================================

	if !security.VerifyOTP(
		req.OTP,
		pending.OTPHash,
	) {
		return nil, apperror.New(
			apperror.CodeInvalidVerifyCode,
			"invalid verification code",
			nil,
		)
	}

	// ============================================================
	// Create User
	// ============================================================

	newUser := &user.User{
		Username:      pending.Username,
		Email:         pending.Email,
		PasswordHash:  pending.PasswordHash,
		Role:          "user",
		Status:        "active",
		EmailVerified: true,
	}

	if err := s.userRepo.Create(
		ctx,
		newUser,
	); err != nil {
		return nil, err
	}

	// ============================================================
	// Delete pending registration
	// ============================================================

	if err := s.redisStore.Delete(
		ctx,
		key,
	); err != nil {

		// User was already created successfully.
		// Don't report registration failure because Redis cleanup
		// failed after the database transaction succeeded.

		return &VerifyRegisterResponse{
			Message: "registration successful",
		}, nil
	}

	// ============================================================
	// Success
	// ============================================================

	return &VerifyRegisterResponse{
		Message: "registration successful",
	}, nil

}
