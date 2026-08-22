package user

import (
	"time"

	dbutils "github.com/khaingminhtun/production-go-api/internal/shared/dbutils"
)

// User has soft delete (DeletedAt) → embed SoftDeleteModel.
type User struct {
	dbutils.SoftDeleteModel

	Email        string `gorm:"size:255;uniqueIndex;not null"`
	Username     string `gorm:"size:100;not null"`
	PasswordHash string `gorm:"not null"`

	Role          Role   `gorm:"size:20;not null;default:'user'"`
	Status        Status `gorm:"size:20;not null;default:'active'"`
	EmailVerified bool   `gorm:"not null;default:false"`

	LastLoginAt *time.Time
}
