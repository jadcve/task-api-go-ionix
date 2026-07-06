package dto

import "time"

type CreateTaskRequest struct {
	Title       string    `json:"title" binding:"required"`
	Description string    `json:"description"`
	DueDate     time.Time `json:"due_date" binding:"required"`
	AssignedTo  uint      `json:"assigned_to" binding:"required"`
}

type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	DueDate     *time.Time `json:"due_date"`
	AssignedTo  *uint      `json:"assigned_to"`
}

type UpdateTaskStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

type CreateTaskCommentRequest struct {
	Comment string `json:"comment" binding:"required"`
}

type TaskCommentResponse struct {
	ID        uint      `json:"id"`
	TaskID    uint      `json:"task_id"`
	UserID    uint      `json:"user_id"`
	Comment   string    `json:"comment"`
	CreatedAt time.Time `json:"created_at"`
}

type TaskResponse struct {
	ID               uint                  `json:"id"`
	Title            string                `json:"title"`
	Description      string                `json:"description"`
	DueDate          time.Time             `json:"due_date"`
	Status           string                `json:"status"`
	AssignedTo       uint                  `json:"assigned_to"`
	AssignedUserName string                `json:"assigned_user_name"`
	CreatedBy        uint                  `json:"created_by"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
	Comments         []TaskCommentResponse `json:"comments,omitempty"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
