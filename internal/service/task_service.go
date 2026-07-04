package service

import (
	"errors"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id uint) (*domain.Task, error)
	FindAll() ([]domain.Task, error)
	Update(task *domain.Task) error
	Delete(id uint) error
	FindAssignedUser(userID uint) (*domain.User, error)
}

type TaskService struct {
	taskRepository TaskRepository
}

func NewTaskService(taskRepository TaskRepository) *TaskService {
	return &TaskService{taskRepository: taskRepository}
}

func (s *TaskService) CreateTask(createdBy uint, request dto.CreateTaskRequest) (*dto.TaskResponse, error) {
	if !request.DueDate.After(time.Now()) {
		return nil, apperrors.ErrInvalidDueDate
	}

	assignedUser, err := s.taskRepository.FindAssignedUser(request.AssignedTo)
	if err != nil {
		return nil, err
	}
	if assignedUser == nil {
		return nil, apperrors.ErrAssignedUserNotFound
	}
	if assignedUser.Role != enums.UserRoleExecutor {
		return nil, apperrors.ErrInvalidUserRole
	}

	task := &domain.Task{
		Title:       strings.TrimSpace(request.Title),
		Description: strings.TrimSpace(request.Description),
		DueDate:     request.DueDate,
		Status:      enums.TaskStatusAssigned,
		AssignedTo:  request.AssignedTo,
		CreatedBy:   createdBy,
	}

	if err := s.taskRepository.Create(task); err != nil {
		return nil, err
	}

	return mapTaskToResponse(task, assignedUser.Name), nil
}

func (s *TaskService) ListTasks() ([]dto.TaskResponse, error) {
	tasks, err := s.taskRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.TaskResponse, 0, len(tasks))
	for i := range tasks {
		assignedUserName := ""
		user, err := s.taskRepository.FindAssignedUser(tasks[i].AssignedTo)
		if err != nil {
			return nil, err
		}
		if user != nil {
			assignedUserName = user.Name
		}
		responses = append(responses, *mapTaskToResponse(&tasks[i], assignedUserName))
	}

	return responses, nil
}

func (s *TaskService) GetTask(id uint) (*dto.TaskResponse, error) {
	task, err := s.taskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperrors.ErrTaskNotFound
	}

	assignedUserName := ""
	user, err := s.taskRepository.FindAssignedUser(task.AssignedTo)
	if err != nil {
		return nil, err
	}
	if user != nil {
		assignedUserName = user.Name
	}

	return mapTaskToResponse(task, assignedUserName), nil
}

func (s *TaskService) UpdateTask(id uint, request dto.UpdateTaskRequest) (*dto.TaskResponse, error) {
	task, err := s.taskRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if task == nil {
		return nil, apperrors.ErrTaskNotFound
	}
	if task.Status != enums.TaskStatusAssigned {
		return nil, apperrors.ErrInvalidTaskStatus
	}

	if request.Title != nil {
		task.Title = strings.TrimSpace(*request.Title)
	}
	if request.Description != nil {
		task.Description = strings.TrimSpace(*request.Description)
	}
	if request.DueDate != nil {
		if !request.DueDate.After(time.Now()) {
			return nil, apperrors.ErrInvalidDueDate
		}
		task.DueDate = *request.DueDate
	}
	if request.AssignedTo != nil {
		assignedUser, err := s.taskRepository.FindAssignedUser(*request.AssignedTo)
		if err != nil {
			return nil, err
		}
		if assignedUser == nil {
			return nil, apperrors.ErrAssignedUserNotFound
		}
		if assignedUser.Role != enums.UserRoleExecutor {
			return nil, apperrors.ErrInvalidUserRole
		}
		task.AssignedTo = *request.AssignedTo
	}

	if err := s.taskRepository.Update(task); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrTaskNotFound
		}
		return nil, err
	}

	assignedUserName := ""
	user, err := s.taskRepository.FindAssignedUser(task.AssignedTo)
	if err != nil {
		return nil, err
	}
	if user != nil {
		assignedUserName = user.Name
	}

	return mapTaskToResponse(task, assignedUserName), nil
}

func (s *TaskService) DeleteTask(id uint) error {
	task, err := s.taskRepository.FindByID(id)
	if err != nil {
		return err
	}
	if task == nil {
		return apperrors.ErrTaskNotFound
	}
	if task.Status != enums.TaskStatusAssigned {
		return apperrors.ErrInvalidTaskStatus
	}

	if err := s.taskRepository.Delete(id); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrTaskNotFound
		}
		return err
	}

	return nil
}

func mapTaskToResponse(task *domain.Task, assignedUserName string) *dto.TaskResponse {
	return &dto.TaskResponse{
		ID:               task.ID,
		Title:            task.Title,
		Description:      task.Description,
		DueDate:          task.DueDate,
		Status:           string(task.Status),
		AssignedTo:       task.AssignedTo,
		AssignedUserName: assignedUserName,
		CreatedBy:        task.CreatedBy,
		CreatedAt:        task.CreatedAt,
		UpdatedAt:        task.UpdatedAt,
	}
}
