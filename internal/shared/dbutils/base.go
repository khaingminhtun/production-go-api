package dbutils

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SoftDeleteModel struct {
	BaseModel

	DeletedAt gorm.DeletedAt `gorm:"index"`
}
