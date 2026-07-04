package service

import (
	"errors"
	"testing"

	"github.com/jackc/pgx/v5"

	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
	"task-api-go-ionix/internal/security"
)

type fakeUserRepository struct {
	usersByEmail map[string]*domain.User
	usersByID    map[uint]*domain.User

	updatedUserID             uint
	updatedPasswordHash       string
	updatedMustChangePassword bool
	updateErr                 error
}

func (r *fakeUserRepository) FindByEmail(email string) (*domain.User, error) {
	if user, ok := r.usersByEmail[email]; ok {
		return user, nil
	}
	return nil, nil
}

func (r *fakeUserRepository) FindByID(id uint) (*domain.User, error) {
	if user, ok := r.usersByID[id]; ok {
		return user, nil
	}
	return nil, nil
}

func (r *fakeUserRepository) UpdatePassword(userID uint, passwordHash string, mustChangePassword bool) error {
	if r.updateErr != nil {
		return r.updateErr
	}

	user, ok := r.usersByID[userID]
	if !ok {
		return pgx.ErrNoRows
	}

	user.PasswordHash = passwordHash
	user.MustChangePassword = mustChangePassword

	r.updatedUserID = userID
	r.updatedPasswordHash = passwordHash
	r.updatedMustChangePassword = mustChangePassword

	return nil
}

func newActiveUser(t *testing.T, id uint, email string, password string) *domain.User {
	t.Helper()

	hash, err := security.HashPassword(password)
	if err != nil {
		t.Fatalf("failed to hash password: %v", err)
	}

	return &domain.User{
		ID:                 id,
		Name:               "Administrator",
		Email:              email,
		PasswordHash:       hash,
		Role:               domain.UserRoleAdmin,
		MustChangePassword: false,
		IsActive:           true,
	}
}

func TestAuthServiceLoginSuccess(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	repo := &fakeUserRepository{
		usersByEmail: map[string]*domain.User{user.Email: user},
		usersByID:    map[uint]*domain.User{user.ID: user},
	}

	svc := NewAuthService(repo, "test-secret", 24)

	res, err := svc.Login(dto.LoginRequest{Email: "admin@test.com", Password: "Admin123"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.Token == "" {
		t.Fatal("expected token to be generated")
	}
	if res.User.Email != "admin@test.com" {
		t.Fatalf("unexpected user email %q", res.User.Email)
	}
}

func TestAuthServiceLoginWrongPassword(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{user.Email: user}, usersByID: map[uint]*domain.User{user.ID: user}}
	svc := NewAuthService(repo, "test-secret", 24)

	_, err := svc.Login(dto.LoginRequest{Email: "admin@test.com", Password: "WrongPassword"})
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceLoginUserNotFound(t *testing.T) {
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{}, usersByID: map[uint]*domain.User{}}
	svc := NewAuthService(repo, "test-secret", 24)

	_, err := svc.Login(dto.LoginRequest{Email: "missing@test.com", Password: "Admin123"})
	if !errors.Is(err, apperrors.ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestAuthServiceLoginInactiveUser(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	user.IsActive = false
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{user.Email: user}, usersByID: map[uint]*domain.User{user.ID: user}}
	svc := NewAuthService(repo, "test-secret", 24)

	_, err := svc.Login(dto.LoginRequest{Email: "admin@test.com", Password: "Admin123"})
	if !errors.Is(err, apperrors.ErrInactiveUser) {
		t.Fatalf("expected ErrInactiveUser, got %v", err)
	}
}

func TestAuthServiceChangePasswordSuccess(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{user.Email: user}, usersByID: map[uint]*domain.User{user.ID: user}}
	svc := NewAuthService(repo, "test-secret", 24)

	err := svc.ChangePassword(1, dto.ChangePasswordRequest{CurrentPassword: "Admin123", NewPassword: "Admin456"})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if repo.updatedUserID != 1 {
		t.Fatalf("expected updated user id 1, got %d", repo.updatedUserID)
	}
	if repo.updatedMustChangePassword {
		t.Fatal("expected must_change_password to be false")
	}
	if !security.CheckPassword("Admin456", repo.updatedPasswordHash) {
		t.Fatal("expected updated password hash to match new password")
	}
}

func TestAuthServiceChangePasswordInvalidCurrentPassword(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{user.Email: user}, usersByID: map[uint]*domain.User{user.ID: user}}
	svc := NewAuthService(repo, "test-secret", 24)

	err := svc.ChangePassword(1, dto.ChangePasswordRequest{CurrentPassword: "WrongCurrent", NewPassword: "Admin456"})
	if !errors.Is(err, apperrors.ErrInvalidCurrentPass) {
		t.Fatalf("expected ErrInvalidCurrentPass, got %v", err)
	}
}

func TestAuthServiceMeSuccess(t *testing.T) {
	user := newActiveUser(t, 1, "admin@test.com", "Admin123")
	repo := &fakeUserRepository{usersByEmail: map[string]*domain.User{user.Email: user}, usersByID: map[uint]*domain.User{user.ID: user}}
	svc := NewAuthService(repo, "test-secret", 24)

	res, err := svc.Me(1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if res.ID != 1 || res.Email != "admin@test.com" {
		t.Fatalf("unexpected response: %+v", res)
	}
}
