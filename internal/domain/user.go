package domain

import (
	"time"

	"task-api-go-ionix/internal/common/enums"
)

type User struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Name               string         `gorm:"size:120;not null" json:"name"`
	Email              string         `gorm:"size:255;not null;uniqueIndex" json:"email"`
	PasswordHash       string         `gorm:"column:password_hash;size:255;not null" json:"-"`
	Role               enums.UserRole `gorm:"size:20;not null" json:"role"`
	MustChangePassword bool           `gorm:"column:must_change_password;not null;default:true" json:"must_change_password"`
	IsActive           bool           `gorm:"column:is_active;not null;default:true" json:"is_active"`
	CreatedAt          time.Time      `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt          time.Time      `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
}
