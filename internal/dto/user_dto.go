package dto

import "time"

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
	Role  string `json:"role" binding:"required"`
}

type UpdateUserRequest struct {
	Name     *string `json:"name"`
	Role     *string `json:"role"`
	IsActive *bool   `json:"is_active"`
}

type UserResponse struct {
	ID                 uint      `json:"id"`
	Name               string    `json:"name"`
	Email              string    `json:"email"`
	Role               string    `json:"role"`
	MustChangePassword bool      `json:"must_change_password"`
	IsActive           bool      `json:"is_active"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

type UserListResponse struct {
	Users []UserResponse `json:"users"`
}
