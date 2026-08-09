package user

import (
	"context"
	"fmt"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, user *User) error
	GetByID(ctx context.Context, id uint) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByUsername(ctx context.Context, username string) (*User, error)
	List(ctx context.Context, offset, limit int) ([]User, int64, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		db: db,
	}
}

func (r *repository) Create(
	ctx context.Context,
	user *User,
) error {
	if err := r.db.WithContext(ctx).
		Create(user).
		Error; err != nil {

		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (r *repository) GetByID(
	ctx context.Context,
	id uint,
) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).
		First(&user, id).
		Error; err != nil {

		return nil, fmt.Errorf(
			"get user by id %d: %w",
			id,
			err,
		)
	}

	return &user, nil
}

func (r *repository) GetByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).
		Where("email = ?", email).
		First(&user).
		Error; err != nil {

		return nil, fmt.Errorf(
			"get user by email: %w",
			err,
		)
	}

	return &user, nil
}

func (r *repository) GetByUsername(
	ctx context.Context,
	username string,
) (*User, error) {
	var user User

	if err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&user).
		Error; err != nil {

		return nil, fmt.Errorf(
			"get user by username: %w",
			err,
		)
	}

	return &user, nil
}

func (r *repository) List(
	ctx context.Context,
	offset, limit int,
) ([]User, int64, error) {
	var users []User
	var total int64

	query := r.db.WithContext(ctx).
		Model(&User{})

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf(
			"count users: %w",
			err,
		)
	}

	if err := query.
		Offset(offset).
		Limit(limit).
		Order("created_at DESC").
		Find(&users).
		Error; err != nil {

		return nil, 0, fmt.Errorf(
			"list users: %w",
			err,
		)
	}

	return users, total, nil
}

func (r *repository) Update(
	ctx context.Context,
	user *User,
) error {
	if err := r.db.WithContext(ctx).
		Save(user).
		Error; err != nil {

		return fmt.Errorf(
			"update user %d: %w",
			user.ID,
			err,
		)
	}

	return nil
}

func (r *repository) Delete(
	ctx context.Context,
	id uint,
) error {
	if err := r.db.WithContext(ctx).
		Delete(&User{}, id).
		Error; err != nil {

		return fmt.Errorf(
			"delete user %d: %w",
			id,
			err,
		)
	}

	return nil
}
