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

type TaskResponse struct {
	ID               uint      `json:"id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	DueDate          time.Time `json:"due_date"`
	Status           string    `json:"status"`
	AssignedTo       uint      `json:"assigned_to"`
	AssignedUserName string    `json:"assigned_user_name"`
	CreatedBy        uint      `json:"created_by"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

type TaskListResponse struct {
	Tasks []TaskResponse `json:"tasks"`
}
