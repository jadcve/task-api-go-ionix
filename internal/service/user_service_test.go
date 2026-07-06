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

type fakeUserManagementRepository struct {
	usersByID    map[uint]*domain.User
	usersByEmail map[string]*domain.User
	nextID       uint
}

func newFakeUserManagementRepository() *fakeUserManagementRepository {
	return &fakeUserManagementRepository{
		usersByID:    map[uint]*domain.User{},
		usersByEmail: map[string]*domain.User{},
		nextID:       1,
	}
}

func (r *fakeUserManagementRepository) Create(user *domain.User) error {
	if _, exists := r.usersByEmail[user.Email]; exists {
		return errors.New("duplicate")
	}

	user.ID = r.nextID
	r.nextID++
	if user.CreatedAt.IsZero() {
		user.CreatedAt = time.Now()
	}
	user.UpdatedAt = user.CreatedAt

	cloned := *user
	r.usersByID[user.ID] = &cloned
	r.usersByEmail[user.Email] = &cloned

	return nil
}

func (r *fakeUserManagementRepository) FindByID(id uint) (*domain.User, error) {
	user, ok := r.usersByID[id]
	if !ok {
		return nil, nil
	}
	cloned := *user
	return &cloned, nil
}

func (r *fakeUserManagementRepository) FindByEmail(email string) (*domain.User, error) {
	user, ok := r.usersByEmail[email]
	if !ok {
		return nil, nil
	}
	cloned := *user
	return &cloned, nil
}

func (r *fakeUserManagementRepository) FindAll() ([]domain.User, error) {
	users := make([]domain.User, 0, len(r.usersByID))
	for _, user := range r.usersByID {
		users = append(users, *user)
	}
	return users, nil
}

func (r *fakeUserManagementRepository) Update(user *domain.User) error {
	existing, ok := r.usersByID[user.ID]
	if !ok {
		return pgx.ErrNoRows
	}

	existing.Name = user.Name
	existing.Role = user.Role
	existing.IsActive = user.IsActive
	existing.UpdatedAt = time.Now()
	user.UpdatedAt = existing.UpdatedAt

	return nil
}

func (r *fakeUserManagementRepository) SoftDelete(id uint) error {
	existing, ok := r.usersByID[id]
	if !ok {
		return pgx.ErrNoRows
	}

	existing.IsActive = false
	existing.UpdatedAt = time.Now()
	return nil
}

func TestCreateUserExecutorSuccess(t *testing.T) {
	repo := newFakeUserManagementRepository()
	svc := NewUserService(repo)

	res, tempPassword, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Juan Ejecutor",
		Email: "juan@test.com",
		Role:  "EXECUTOR",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if tempPassword != "ChangeMe123!" {
		t.Fatalf("expected fixed temporary password, got %s", tempPassword)
	}
	if res.Role != "EXECUTOR" {
		t.Fatalf("expected role EXECUTOR, got %s", res.Role)
	}
}

func TestCreateUserAuditorSuccess(t *testing.T) {
	repo := newFakeUserManagementRepository()
	svc := NewUserService(repo)

	res, _, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Ana Auditor",
		Email: "ana@test.com",
		Role:  "AUDITOR",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Role != "AUDITOR" {
		t.Fatalf("expected role AUDITOR, got %s", res.Role)
	}
}

func TestCreateUserAdminNotAllowed(t *testing.T) {
	repo := newFakeUserManagementRepository()
	svc := NewUserService(repo)

	_, _, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Other Admin",
		Email: "other-admin@test.com",
		Role:  "ADMIN",
	})
	if !errors.Is(err, apperrors.ErrInvalidUserRole) {
		t.Fatalf("expected ErrInvalidUserRole, got %v", err)
	}
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	repo := newFakeUserManagementRepository()
	repo.usersByID[1] = &domain.User{ID: 1, Name: "Already", Email: "dup@test.com", Role: enums.UserRoleExecutor, IsActive: true}
	repo.usersByEmail["dup@test.com"] = repo.usersByID[1]
	repo.nextID = 2

	svc := NewUserService(repo)

	_, _, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Duplicate",
		Email: "dup@test.com",
		Role:  "EXECUTOR",
	})
	if !errors.Is(err, apperrors.ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestCreateUserMustChangePasswordAndActive(t *testing.T) {
	repo := newFakeUserManagementRepository()
	svc := NewUserService(repo)

	res, _, err := svc.CreateUser(dto.CreateUserRequest{
		Name:  "Mario",
		Email: "mario@test.com",
		Role:  "EXECUTOR",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !res.MustChangePassword {
		t.Fatal("expected must_change_password=true")
	}
	if !res.IsActive {
		t.Fatal("expected is_active=true")
	}
}

func TestUpdateUserAdminRoleNotAllowed(t *testing.T) {
	repo := newFakeUserManagementRepository()
	repo.usersByID[1] = &domain.User{ID: 1, Name: "User", Email: "user@test.com", Role: enums.UserRoleExecutor, IsActive: true}
	repo.usersByEmail["user@test.com"] = repo.usersByID[1]
	svc := NewUserService(repo)

	role := "ADMIN"
	_, err := svc.UpdateUser(1, dto.UpdateUserRequest{Role: &role})
	if !errors.Is(err, apperrors.ErrInvalidUserRole) {
		t.Fatalf("expected ErrInvalidUserRole, got %v", err)
	}
}

func TestDeleteUserSoftDelete(t *testing.T) {
	repo := newFakeUserManagementRepository()
	repo.usersByID[1] = &domain.User{ID: 1, Name: "Delete Me", Email: "delete@test.com", Role: enums.UserRoleExecutor, IsActive: true}
	repo.usersByEmail["delete@test.com"] = repo.usersByID[1]
	svc := NewUserService(repo)

	err := svc.DeleteUser(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.usersByID[1].IsActive {
		t.Fatal("expected user to be soft deleted (is_active=false)")
	}
}
