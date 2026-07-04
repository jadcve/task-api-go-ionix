package service

import (
	"errors"
	"testing"
	"time"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
)

type fakeTaskRepository struct {
	tasks  map[uint]*domain.Task
	users  map[uint]*domain.User
	nextID uint
}

func newFakeTaskRepository() *fakeTaskRepository {
	return &fakeTaskRepository{
		tasks:  map[uint]*domain.Task{},
		users:  map[uint]*domain.User{},
		nextID: 1,
	}
}

func (r *fakeTaskRepository) Create(task *domain.Task) error {
	task.ID = r.nextID
	r.nextID++
	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now
	cloned := *task
	r.tasks[task.ID] = &cloned
	return nil
}

func (r *fakeTaskRepository) FindByID(id uint) (*domain.Task, error) {
	task, ok := r.tasks[id]
	if !ok {
		return nil, nil
	}
	cloned := *task
	return &cloned, nil
}

func (r *fakeTaskRepository) FindAll() ([]domain.Task, error) {
	items := make([]domain.Task, 0, len(r.tasks))
	for _, t := range r.tasks {
		if t.DeletedAt == nil {
			items = append(items, *t)
		}
	}
	return items, nil
}

func (r *fakeTaskRepository) Update(task *domain.Task) error {
	existing, ok := r.tasks[task.ID]
	if !ok || existing.DeletedAt != nil {
		return errors.New("not found")
	}
	existing.Title = task.Title
	existing.Description = task.Description
	existing.DueDate = task.DueDate
	existing.AssignedTo = task.AssignedTo
	existing.UpdatedAt = time.Now()
	task.UpdatedAt = existing.UpdatedAt
	return nil
}

func (r *fakeTaskRepository) Delete(id uint) error {
	existing, ok := r.tasks[id]
	if !ok || existing.DeletedAt != nil {
		return errors.New("not found")
	}
	now := time.Now()
	existing.DeletedAt = &now
	existing.UpdatedAt = now
	return nil
}

func (r *fakeTaskRepository) FindAssignedUser(userID uint) (*domain.User, error) {
	user, ok := r.users[userID]
	if !ok {
		return nil, nil
	}
	cloned := *user
	return &cloned, nil
}

func TestCreateTaskSuccessInitialStatusAssigned(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	svc := NewTaskService(repo)

	res, err := svc.CreateTask(1, dto.CreateTaskRequest{
		Title: "Task 1", Description: "Desc", DueDate: time.Now().Add(2 * time.Hour), AssignedTo: 10,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusAssigned) {
		t.Fatalf("expected ASSIGNED, got %s", res.Status)
	}
}

func TestCreateTaskCannotAssignAuditor(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[20] = &domain.User{ID: 20, Name: "Aud", Role: enums.UserRoleAuditor, IsActive: true}
	svc := NewTaskService(repo)

	_, err := svc.CreateTask(1, dto.CreateTaskRequest{Title: "Task", DueDate: time.Now().Add(time.Hour), AssignedTo: 20})
	if !errors.Is(err, apperrors.ErrInvalidUserRole) {
		t.Fatalf("expected ErrInvalidUserRole, got %v", err)
	}
}

func TestCreateTaskCannotAssignAdmin(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[30] = &domain.User{ID: 30, Name: "Admin", Role: enums.UserRoleAdmin, IsActive: true}
	svc := NewTaskService(repo)

	_, err := svc.CreateTask(1, dto.CreateTaskRequest{Title: "Task", DueDate: time.Now().Add(time.Hour), AssignedTo: 30})
	if !errors.Is(err, apperrors.ErrInvalidUserRole) {
		t.Fatalf("expected ErrInvalidUserRole, got %v", err)
	}
}

func TestCreateTaskPastDueDate(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	svc := NewTaskService(repo)

	_, err := svc.CreateTask(1, dto.CreateTaskRequest{Title: "Task", DueDate: time.Now().Add(-time.Hour), AssignedTo: 10})
	if !errors.Is(err, apperrors.ErrInvalidDueDate) {
		t.Fatalf("expected ErrInvalidDueDate, got %v", err)
	}
}

func TestUpdateTaskOnlyAssigned(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusInProgress, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	title := "Updated"
	_, err := svc.UpdateTask(1, dto.UpdateTaskRequest{Title: &title})
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

func TestDeleteTaskOnlyAssigned(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusCompleted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	err := svc.DeleteTask(1)
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}
