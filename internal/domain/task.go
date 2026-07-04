package domain

import (
	"time"

	"task-api-go-ionix/internal/common/enums"
)

type Task struct {
	ID          uint             `gorm:"primaryKey" json:"id"`
	Title       string           `gorm:"size:200;not null" json:"title"`
	Description string           `gorm:"type:text" json:"description"`
	DueDate     time.Time        `gorm:"column:due_date;not null" json:"due_date"`
	Status      enums.TaskStatus `gorm:"size:20;not null" json:"status"`
	AssignedTo  uint             `gorm:"column:assigned_to;not null;index" json:"assigned_to"`
	CreatedBy   uint             `gorm:"column:created_by;not null;index" json:"created_by"`
	CreatedAt   time.Time        `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`
	DeletedAt   *time.Time       `gorm:"column:deleted_at" json:"-"`

	AssignedToUser User `gorm:"foreignKey:AssignedTo;references:ID" json:"-"`
	CreatedByUser  User `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}

func (t Task) IsExpired(now time.Time) bool {
	if t.DueDate.IsZero() {
		return false
	}

	return t.DueDate.Before(now)
}
