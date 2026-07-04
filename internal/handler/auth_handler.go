package handler

import (
	"errors"

	"github.com/gin-gonic/gin"

	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
	"task-api-go-ionix/internal/response"
	"task-api-go-ionix/internal/service"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var request dto.LoginRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	loginResponse, err := h.authService.Login(request)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCredentials):
			response.Error(c, 401, "Invalid credentials", nil)
		case errors.Is(err, apperrors.ErrInactiveUser):
			response.Error(c, 403, "Inactive user", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "Login successful", loginResponse)
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	var request dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	err := h.authService.ChangePassword(userID, request)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrInvalidCurrentPass):
			response.Error(c, 401, "Invalid current password", nil)
		case errors.Is(err, apperrors.ErrUserNotFound):
			response.Error(c, 404, "User not found", nil)
		case errors.Is(err, apperrors.ErrInactiveUser):
			response.Error(c, 403, "Inactive user", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "Password changed successfully", nil)
}

func (h *AuthHandler) Me(c *gin.Context) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	user, err := h.authService.Me(userID)
	if err != nil {
		switch {
		case errors.Is(err, apperrors.ErrUserNotFound):
			response.Error(c, 404, "User not found", nil)
		case errors.Is(err, apperrors.ErrInactiveUser):
			response.Error(c, 403, "Inactive user", nil)
		default:
			response.Error(c, 500, "Internal server error", nil)
		}
		return
	}

	response.Success(c, 200, "Authenticated user retrieved successfully", user)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	response.NoContent(c)
}
