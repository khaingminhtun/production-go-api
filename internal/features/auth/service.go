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
	"github.com/khaingminhtun/production-go-api/internal/shared/security"
)

type Service interface {
	Register(
		ctx context.Context,
		req RegisterRequest) (*RegisterResponse, error)
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
	// registration logic

	username := strings.TrimSpace(req.Username)
	email := strings.ToLower(strings.TrimSpace(req.Email))

	_, err := s.userRepo.GetByEmail(ctx, email)
	if err == nil {
		return nil, err
	}

	_, err = s.userRepo.GetByUsername(ctx, username)
	if err == nil {
		return nil, err
	}

	// hash
	passwordHash, err := security.HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	otp, err := security.GenerateOTP()
	if err != nil {
		return nil, fmt.Errorf("generate otp: %w", err)
	}

	otpHash := security.HashOTP(otp)

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
	// Save registration in Redis
	// ============================================================

	key := "auth:register:" + email

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
	// Publish email job to Redis Stream
	// ============================================================

	if err := s.emailQueue.Publish(
		ctx,
		job,
	); err != nil {

		// Don't leave a registration waiting for an email
		// that could not be queued.
		_ = s.redisStore.Delete(
			ctx,
			key,
		)

		return nil, fmt.Errorf(
			"queue verification email: %w",
			err,
		)
	}

	// ============================================================
	// Response
	// ============================================================

	return &RegisterResponse{
		Message: "Verification code sent to your email",
	}, nil

}
