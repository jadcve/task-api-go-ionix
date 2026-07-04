package domain

import "time"

type TaskComment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TaskID    uint      `gorm:"column:task_id;not null;index" json:"task_id"`
	UserID    uint      `gorm:"column:user_id;not null;index" json:"user_id"`
	Comment   string    `gorm:"type:text;not null" json:"comment"`
	CreatedAt time.Time `gorm:"column:created_at;not null;autoCreateTime" json:"created_at"`

	Task Task `gorm:"foreignKey:TaskID;references:ID" json:"-"`
	User User `gorm:"foreignKey:UserID;references:ID" json:"-"`
}
