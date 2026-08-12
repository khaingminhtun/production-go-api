package user

import (
	"context"
	"errors"
	"strings"

	postgres "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
	"github.com/khaingminhtun/production-go-api/internal/shared/errorHandler/apperror"
	"gorm.io/gorm"
)

type Service interface {
	CreateUser(
		ctx context.Context,
		req CreateUserRequest,
	) (*UserResponse, error)

	GetUser(
		ctx context.Context,
		id uint,
	) (*UserResponse, error)

	ListUsers(
		ctx context.Context,
		offset, limit int,
	) (*UserListResponse, error)

	UpdateUser(
		ctx context.Context,
		id uint,
		req UpdateUserRequest,
	) (*UserResponse, error)

	DeleteUser(
		ctx context.Context,
		id uint,
	) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

// ============================================================
// Create
// ============================================================

func (s *service) CreateUser(
	ctx context.Context,
	req CreateUserRequest,
) (*UserResponse, error) {

	email := strings.ToLower(strings.TrimSpace(req.Email))
	username := strings.TrimSpace(req.Username)

	// Application-level check.
	_, err := s.repo.GetByEmail(ctx, email)

	if err == nil {
		return nil, apperror.New(
			apperror.CodeUserAlreadyExists,
			"user with this email already exists",
			nil,
		)
	}

	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	//passwordHash, err := password.Hash(req.Password)
	//if err != nil {
	//	return nil, fmt.Errorf("hash password: %w", err)
	//}

	user := &User{
		Email:         email,
		Username:      username,
		PasswordHash:  req.Password,
		Role:          "user",
		Status:        "active",
		EmailVerified: false,
	}

	if err := s.repo.Create(ctx, user); err != nil {

		// PostgreSQL unique constraint is the final authority.
		if postgres.IsUniqueViolation(err) &&
			postgres.ConstraintName(err) == "users_email_key" {

			return nil, apperror.New(
				apperror.CodeUserAlreadyExists,
				"user with this email already exists",
				err,
			)
		}

		return nil, err
	}

	return toUserResponse(user), nil
}

// ============================================================
// Get
// ============================================================

func (s *service) GetUser(
	ctx context.Context,
	id uint,
) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userNotFoundError(err)
		}

		return nil, err
	}

	return toUserResponse(user), nil
}

// ============================================================
// List
// ============================================================

func (s *service) ListUsers(
	ctx context.Context,
	offset, limit int,
) (*UserListResponse, error) {

	if offset < 0 {
		offset = 0
	}

	if limit <= 0 {
		limit = 20
	}

	if limit > 100 {
		limit = 100
	}

	users, total, err := s.repo.List(
		ctx,
		offset,
		limit,
	)
	if err != nil {
		return nil, err
	}

	return &UserListResponse{
		Users:  toUserResponseList(users),
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

// ============================================================
// Update
// ============================================================

func (s *service) UpdateUser(
	ctx context.Context,
	id uint,
	req UpdateUserRequest,
) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, userNotFoundError(err)
		}

		return nil, err
	}

	if req.Email != "" {
		user.Email = strings.ToLower(
			strings.TrimSpace(req.Email),
		)
	}

	if req.Username != "" {
		user.Username = strings.TrimSpace(req.Username)
	}

	if err := s.repo.Update(ctx, user); err != nil {

		if postgres.IsUniqueViolation(err) &&
			postgres.ConstraintName(err) == "users_email_key" {

			return nil, apperror.New(
				apperror.CodeUserAlreadyExists,
				"user with this email already exists",
				err,
			)
		}

		return nil, err
	}

	return toUserResponse(user), nil
}

// ============================================================
// Delete
// ============================================================

func (s *service) DeleteUser(
	ctx context.Context,
	id uint,
) error {

	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userNotFoundError(err)
		}

		return err
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	return nil
}

// ============================================================
// Error Helpers
// ============================================================

func userNotFoundError(err error) error {
	return apperror.New(
		apperror.CodeUserNotFound,
		"user not found",
		err,
	)
}
