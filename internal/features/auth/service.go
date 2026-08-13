package auth

import (
	"context"
	"fmt"
	"strings"

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
	userRepo       user.Repository
	redisStore     redisinfra.RedisStore
	redisPublisher redisinfra.RedisPublisher
}

func NewService(
	userRep user.Repository,
	redisStore redisinfra.RedisStore,
	redisPublisher redisinfra.RedisPublisher,
) Service {
	return &service{
		userRepo:       userRep,
		redisStore:     redisStore,
		redisPublisher: redisPublisher,
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

}
