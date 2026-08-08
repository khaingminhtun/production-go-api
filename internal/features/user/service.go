package user

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrEmailAlreadyUsed = errors.New("email already in use")
)

type Service interface {
	// Normal authenticated user.
	GetMe(ctx context.Context, userID uint) (*UserResponse, error)

	UpdateEmail(
		ctx context.Context,
		userID uint,
		req UpdateEmailRequest,
	) (*UserResponse, error)

	// Admin.
	ListUsers(
		ctx context.Context,
		offset, limit int,
	) (*UserListResponse, error)

	GetUser(
		ctx context.Context,
		id uint,
	) (*AdminUserResponse, error)

	UpdateUser(
		ctx context.Context,
		id uint,
		req AdminUpdateUserRequest,
	) (*AdminUserResponse, error)

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

func (s *service) GetMe(ctx context.Context, userID uint) (*UserResponse, error) {

	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return toUserResponse(user), nil
}

func (s *service) UpdateEmail(
	ctx context.Context,
	userID uint,
	req UpdateEmailRequest,
) (*UserResponse, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Nothing changed.
	if email == user.Email {
		return toUserResponse(user), nil
	}

	// Email must be unique.
	existing, err := s.repo.GetByEmail(ctx, email)

	if err == nil && existing.ID != user.ID {
		return nil, ErrEmailAlreadyUsed
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	user.Email = email

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toUserResponse(user), nil
}

// admin task
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

	users, total, err := s.repo.List(ctx, offset, limit)
	if err != nil {
		return nil, err
	}

	result := make([]AdminUserResponse, 0, len(users))

	for i := range users {
		result = append(
			result,
			*toAdminUserResponse(&users[i]),
		)
	}

	return &UserListResponse{
		Users:  result,
		Total:  total,
		Offset: offset,
		Limit:  limit,
	}, nil
}

func (s *service) GetUser(
	ctx context.Context,
	id uint,
) (*AdminUserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	return toAdminUserResponse(user), nil
}

func (s *service) UpdateUser(
	ctx context.Context,
	id uint,
	req AdminUpdateUserRequest,
) (*AdminUserResponse, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}

		return nil, err
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	// Check email uniqueness only if it changed.
	if email != user.Email {
		existing, err := s.repo.GetByEmail(ctx, email)

		if err == nil && existing.ID != user.ID {
			return nil, ErrEmailAlreadyUsed
		}

		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	user.Email = email

	if req.Role != "" {
		user.Role = req.Role
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}

	return toAdminUserResponse(user), nil
}

func (s *service) DeleteUser(
	ctx context.Context,
	id uint,
) error {
	_, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}

		return err
	}

	return s.repo.Delete(ctx, id)
}
