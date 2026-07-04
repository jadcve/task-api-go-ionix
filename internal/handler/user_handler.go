package handler

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
	"task-api-go-ionix/internal/response"
	"task-api-go-ionix/internal/service"
)

type UserHandler struct {
	userService *service.UserService
}

func NewUserHandler(userService *service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var request dto.CreateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	user, temporaryPassword, err := h.userService.CreateUser(request)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrEmailAlreadyExists):
			response.Error(c, 409, "Email already exists", nil)
		case errors.Is(err, apperrors.ErrInvalidUserRole):
			response.Error(c, 400, "Invalid user role", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 201, "User created successfully", gin.H{
		"user":               user,
		"temporary_password": temporaryPassword,
	})
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	users, err := h.userService.ListUsers()
	if err != nil {
		response.Error(c, 500, "Internal server error", nil)
		return
	}

	response.Success(c, 200, "Users retrieved successfully", users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Error(c, 400, "Invalid user id", nil)
		return
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			response.Error(c, 404, "User not found", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "User retrieved successfully", user)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Error(c, 400, "Invalid user id", nil)
		return
	}

	var request dto.UpdateUserRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	user, err := h.userService.UpdateUser(id, request)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidUserRole):
			response.Error(c, 400, "Invalid user role", nil)
		case errors.Is(err, apperrors.ErrUserNotFound):
			response.Error(c, 404, "User not found", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "User updated successfully", user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		response.Error(c, 400, "Invalid user id", nil)
		return
	}

	err := h.userService.DeleteUser(id)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			response.Error(c, 404, "User not found", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "User deleted successfully", nil)
}

func parseID(c *gin.Context) (uint, bool) {
	idRaw := c.Param("id")
	parsed, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(parsed), true
}
