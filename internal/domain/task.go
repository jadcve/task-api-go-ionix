package domain

import "time"

type TaskStatus string

const (
	TaskStatusAssigned   TaskStatus = "ASSIGNED"
	TaskStatusInProgress TaskStatus = "IN_PROGRESS"
	TaskStatusCompleted  TaskStatus = "COMPLETED"
	TaskStatusCancelled  TaskStatus = "CANCELLED"
)

type Task struct {
	ID          uint       `gorm:"primaryKey" json:"id"`
	Title       string     `gorm:"size:200;not null" json:"title"`
	Description string     `gorm:"type:text" json:"description"`
	DueDate     *time.Time `gorm:"column:due_date" json:"due_date"`
	Status      TaskStatus `gorm:"size:20;not null" json:"status"`
	AssignedTo  uint       `gorm:"column:assigned_to;not null;index" json:"assigned_to"`
	CreatedBy   uint       `gorm:"column:created_by;not null;index" json:"created_by"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null;autoUpdateTime" json:"updated_at"`

	AssignedToUser User `gorm:"foreignKey:AssignedTo;references:ID" json:"-"`
	CreatedByUser  User `gorm:"foreignKey:CreatedBy;references:ID" json:"-"`
}
