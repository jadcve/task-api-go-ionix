package service

import (
	"errors"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
)

type fakeTaskRepository struct {
	tasks            map[uint]*domain.Task
	commentsByTaskID map[uint][]domain.TaskComment
	users            map[uint]*domain.User
	nextID           uint
	nextCommentID    uint
}

func newFakeTaskRepository() *fakeTaskRepository {
	return &fakeTaskRepository{
		tasks:            map[uint]*domain.Task{},
		commentsByTaskID: map[uint][]domain.TaskComment{},
		users:            map[uint]*domain.User{},
		nextID:           1,
		nextCommentID:    1,
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
	if !ok || task.DeletedAt != nil {
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
		return pgx.ErrNoRows
	}
	existing.Title = task.Title
	existing.Description = task.Description
	existing.DueDate = task.DueDate
	existing.Status = task.Status
	existing.AssignedTo = task.AssignedTo
	existing.UpdatedAt = time.Now()
	task.UpdatedAt = existing.UpdatedAt
	return nil
}

func (r *fakeTaskRepository) Delete(id uint) error {
	existing, ok := r.tasks[id]
	if !ok || existing.DeletedAt != nil {
		return pgx.ErrNoRows
	}
	now := time.Now()
	existing.DeletedAt = &now
	existing.UpdatedAt = now
	return nil
}

func (r *fakeTaskRepository) FindAllByAssignedTo(userID uint) ([]domain.Task, error) {
	items := make([]domain.Task, 0)
	for _, t := range r.tasks {
		if t.AssignedTo == userID && t.DeletedAt == nil {
			items = append(items, *t)
		}
	}
	return items, nil
}

func (r *fakeTaskRepository) FindAssignedUser(userID uint) (*domain.User, error) {
	user, ok := r.users[userID]
	if !ok {
		return nil, nil
	}
	cloned := *user
	return &cloned, nil
}

func (r *fakeTaskRepository) AddComment(comment *domain.TaskComment) error {
	comment.ID = r.nextCommentID
	r.nextCommentID++
	comment.CreatedAt = time.Now()

	cloned := *comment
	r.commentsByTaskID[comment.TaskID] = append(r.commentsByTaskID[comment.TaskID], cloned)
	return nil
}

func (r *fakeTaskRepository) FindCommentsByTaskID(taskID uint) ([]domain.TaskComment, error) {
	source := r.commentsByTaskID[taskID]
	if len(source) == 0 {
		return []domain.TaskComment{}, nil
	}

	items := make([]domain.TaskComment, 0, len(source))
	for i := range source {
		items = append(items, source[i])
	}
	return items, nil
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
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusStarted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	title := "Updated"
	_, err := svc.UpdateTask(1, dto.UpdateTaskRequest{Title: &title})
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

func TestDeleteTaskOnlyAssigned(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusCompletedSuccess, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	err := svc.DeleteTask(1)
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

func TestListMyTasksReturnsOnlyExecutorTasks(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec 10", Role: enums.UserRoleExecutor, IsActive: true}
	repo.users[20] = &domain.User{ID: 20, Name: "Exec 20", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Mine", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	repo.tasks[2] = &domain.Task{ID: 2, Title: "Other", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 20, CreatedBy: 1}
	svc := NewTaskService(repo)

	items, err := svc.ListMyTasks(10)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 task, got %d", len(items))
	}
	if items[0].AssignedTo != 10 {
		t.Fatalf("expected assigned_to=10, got %d", items[0].AssignedTo)
	}
}

func TestGetMyTaskAllowsOwnedTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Mine", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	task, err := svc.GetMyTask(10, 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if task == nil || task.ID != 1 {
		t.Fatalf("expected task id 1")
	}
}

func TestGetMyTaskBlocksForeignTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Not mine", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 20, CreatedBy: 1}
	svc := NewTaskService(repo)

	_, err := svc.GetMyTask(10, 1)
	if !errors.Is(err, apperrors.ErrForbiddenTaskAccess) {
		t.Fatalf("expected ErrForbiddenTaskAccess, got %v", err)
	}
}

func TestUpdateMyTaskStatusAllowsAssignedToStarted(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	res, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "STARTED"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusStarted) {
		t.Fatalf("expected STARTED, got %s", res.Status)
	}
}

func TestUpdateMyTaskStatusAllowsStartedToWaiting(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusStarted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	res, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "WAITING"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusWaiting) {
		t.Fatalf("expected WAITING, got %s", res.Status)
	}
}

func TestUpdateMyTaskStatusAllowsWaitingToStarted(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusWaiting, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	res, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "STARTED"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusStarted) {
		t.Fatalf("expected STARTED, got %s", res.Status)
	}
}

func TestUpdateMyTaskStatusAllowsStartedToCompletedSuccess(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusStarted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	res, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "COMPLETED_SUCCESS"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusCompletedSuccess) {
		t.Fatalf("expected COMPLETED_SUCCESS, got %s", res.Status)
	}
}

func TestUpdateMyTaskStatusAllowsStartedToCompletedError(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusStarted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	res, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "COMPLETED_ERROR"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Status != string(enums.TaskStatusCompletedError) {
		t.Fatalf("expected COMPLETED_ERROR, got %s", res.Status)
	}
}

func TestUpdateMyTaskStatusBlocksInvalidTransition(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	_, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "WAITING"})
	if !errors.Is(err, apperrors.ErrInvalidTaskStatusTransition) {
		t.Fatalf("expected ErrInvalidTaskStatusTransition, got %v", err)
	}
}

func TestUpdateMyTaskStatusBlocksExpiredTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Expired", DueDate: time.Now().Add(-time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	_, err := svc.UpdateMyTaskStatus(10, 1, dto.UpdateTaskStatusRequest{Status: "STARTED"})
	if !errors.Is(err, apperrors.ErrTaskExpired) {
		t.Fatalf("expected ErrTaskExpired, got %v", err)
	}
}

func TestAdminCannotUpdateTaskInStarted(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusStarted, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	title := "Blocked"
	_, err := svc.UpdateTask(1, dto.UpdateTaskRequest{Title: &title})
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

func TestAdminCannotDeleteTaskInCompletedSuccess(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Task", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusCompletedSuccess, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	err := svc.DeleteTask(1)
	if !errors.Is(err, apperrors.ErrInvalidTaskStatus) {
		t.Fatalf("expected ErrInvalidTaskStatus, got %v", err)
	}
}

func TestAddExpiredTaskCommentAllowsExpiredTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Expired", DueDate: time.Now().Add(-time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	comment, err := svc.AddExpiredTaskComment(10, 1, dto.CreateTaskCommentRequest{Comment: "Bloqueada por vencimiento"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if comment == nil || comment.TaskID != 1 || comment.UserID != 10 {
		t.Fatalf("expected stored comment for task=1 user=10")
	}
}

func TestAddExpiredTaskCommentBlocksNonExpiredTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Future", DueDate: time.Now().Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	svc := NewTaskService(repo)

	_, err := svc.AddExpiredTaskComment(10, 1, dto.CreateTaskCommentRequest{Comment: "No debe pasar"})
	if !errors.Is(err, apperrors.ErrTaskNotExpired) {
		t.Fatalf("expected ErrTaskNotExpired, got %v", err)
	}
}

func TestAddExpiredTaskCommentBlocksForeignTask(t *testing.T) {
	repo := newFakeTaskRepository()
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Expired", DueDate: time.Now().Add(-time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 20, CreatedBy: 1}
	svc := NewTaskService(repo)

	_, err := svc.AddExpiredTaskComment(10, 1, dto.CreateTaskCommentRequest{Comment: "No debe pasar"})
	if !errors.Is(err, apperrors.ErrForbiddenTaskAccess) {
		t.Fatalf("expected ErrForbiddenTaskAccess, got %v", err)
	}
}

func TestListTasksForAuditorReturnsNonDeletedTasks(t *testing.T) {
	repo := newFakeTaskRepository()
	now := time.Now()
	deletedAt := now
	repo.users[10] = &domain.User{ID: 10, Name: "Exec", Role: enums.UserRoleExecutor, IsActive: true}
	repo.tasks[1] = &domain.Task{ID: 1, Title: "Visible", DueDate: now.Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1}
	repo.tasks[2] = &domain.Task{ID: 2, Title: "Deleted", DueDate: now.Add(time.Hour), Status: enums.TaskStatusAssigned, AssignedTo: 10, CreatedBy: 1, DeletedAt: &deletedAt}
	svc := NewTaskService(repo)

	items, err := svc.ListTasksForAuditor()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(items) != 1 {
		t.Fatalf("expected 1 visible task, got %d", len(items))
	}
	if items[0].ID != 1 {
		t.Fatalf("expected task id 1, got %d", items[0].ID)
	}
}
