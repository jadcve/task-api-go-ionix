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

type TaskHandler struct {
	taskService *service.TaskService
}

func NewTaskHandler(taskService *service.TaskService) *TaskHandler {
	return &TaskHandler{taskService: taskService}
}

func (h *TaskHandler) CreateTask(c *gin.Context) {
	userID, ok := getUserID(c)
	if !ok {
		response.Error(c, 401, "Unauthorized", nil)
		return
	}

	var request dto.CreateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	task, err := h.taskService.CreateTask(userID, request)
	if err != nil {
		handleTaskServiceError(c, err)
		return
	}

	response.Success(c, 201, "Task created successfully", task)
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.taskService.ListTasks()
	if err != nil {
		response.Error(c, 500, "Internal server error", nil)
		return
	}

	response.Success(c, 200, "Tasks retrieved successfully", tasks)
}

func (h *TaskHandler) GetTaskByID(c *gin.Context) {
	id, ok := parseTaskID(c)
	if !ok {
		response.Error(c, 400, "Invalid task id", nil)
		return
	}

	task, err := h.taskService.GetTask(id)
	if err != nil {
		handleTaskServiceError(c, err)
		return
	}

	response.Success(c, 200, "Task retrieved successfully", task)
}

func (h *TaskHandler) UpdateTask(c *gin.Context) {
	id, ok := parseTaskID(c)
	if !ok {
		response.Error(c, 400, "Invalid task id", nil)
		return
	}

	var request dto.UpdateTaskRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		response.Error(c, 400, "Invalid request", nil)
		return
	}

	task, err := h.taskService.UpdateTask(id, request)
	if err != nil {
		handleTaskServiceError(c, err)
		return
	}

	response.Success(c, 200, "Task updated successfully", task)
}

func (h *TaskHandler) DeleteTask(c *gin.Context) {
	id, ok := parseTaskID(c)
	if !ok {
		response.Error(c, 400, "Invalid task id", nil)
		return
	}

	err := h.taskService.DeleteTask(id)
	if err != nil {
		handleTaskServiceError(c, err)
		return
	}

	response.Success(c, 200, "Task deleted successfully", nil)
}

func parseTaskID(c *gin.Context) (uint, bool) {
	idRaw := c.Param("id")
	parsed, err := strconv.ParseUint(idRaw, 10, 64)
	if err != nil {
		return 0, false
	}

	return uint(parsed), true
}

func getUserID(c *gin.Context) (uint, bool) {
	userIDValue, exists := c.Get("user_id")
	if !exists {
		return 0, false
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		return 0, false
	}

	return userID, true
}

func handleTaskServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, apperrors.ErrTaskNotFound):
		response.Error(c, 404, "Task not found", nil)
	case errors.Is(err, apperrors.ErrAssignedUserNotFound):
		response.Error(c, 409, "Assigned user not found", nil)
	case errors.Is(err, apperrors.ErrInvalidUserRole):
		response.Error(c, 400, "Invalid user role", nil)
	case errors.Is(err, apperrors.ErrInvalidTaskStatus):
		response.Error(c, 400, "Invalid task status", nil)
	case errors.Is(err, apperrors.ErrInvalidDueDate):
		response.Error(c, 400, "Invalid due date", nil)
	default:
		response.Error(c, 500, "Internal server error", nil)
	}
}
