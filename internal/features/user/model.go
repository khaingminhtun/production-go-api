package user

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint   `gorm:"primaryKey"`
	Email        string `gorm:"size:255;uniqueIndex;not null"`
	Username     string `gorm:"size:100;uniqueIndex;not null"`
	PasswordHash string `gorm:"not null"`

	Role          string `gorm:"size:20;not null;default:'user'"`
	Status        string `gorm:"size:20;not null;default:'active'"`
	EmailVerified bool   `gorm:"not null;default:false"`

	LastLoginAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}
