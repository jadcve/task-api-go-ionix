package service

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	"task-api-go-ionix/internal/common/enums"
	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
	"task-api-go-ionix/internal/security"
)

type UserManagementRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindAll() ([]domain.User, error)
	Update(user *domain.User) error
	SoftDelete(id uint) error
}

type UserService struct {
	userRepository UserManagementRepository
}

func NewUserService(userRepository UserManagementRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) CreateUser(request dto.CreateUserRequest) (*dto.UserResponse, string, error) {
	role := strings.ToUpper(strings.TrimSpace(request.Role))
	if role != string(enums.UserRoleExecutor) && role != string(enums.UserRoleAuditor) {
		return nil, "", apperrors.ErrInvalidUserRole
	}

	existingUser, err := s.userRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, "", err
	}
	if existingUser != nil {
		return nil, "", apperrors.ErrEmailAlreadyExists
	}

	temporaryPassword, err := generateTemporaryPassword()
	if err != nil {
		return nil, "", err
	}

	hash, err := security.HashPassword(temporaryPassword)
	if err != nil {
		return nil, "", err
	}

	user := &domain.User{
		Name:               request.Name,
		Email:              request.Email,
		PasswordHash:       hash,
		Role:               enums.UserRole(role),
		MustChangePassword: true,
		IsActive:           true,
	}

	err = s.userRepository.Create(user)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, "", apperrors.ErrEmailAlreadyExists
		}
		return nil, "", err
	}

	// For this technical challenge we return the temporary password in the API response.
	// In production this should be delivered through a secure channel such as email.
	return mapUserToResponse(user), temporaryPassword, nil
}

func (s *UserService) GetUserByID(id uint) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	return mapUserToResponse(user), nil
}

func (s *UserService) ListUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepository.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]dto.UserResponse, 0, len(users))
	for i := range users {
		responses = append(responses, *mapUserToResponse(&users[i]))
	}

	return responses, nil
}

func (s *UserService) UpdateUser(id uint, request dto.UpdateUserRequest) (*dto.UserResponse, error) {
	user, err := s.userRepository.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}

	if request.Name != nil {
		user.Name = *request.Name
	}

	if request.Role != nil {
		role := strings.ToUpper(strings.TrimSpace(*request.Role))
		if role != string(enums.UserRoleExecutor) && role != string(enums.UserRoleAuditor) {
			return nil, apperrors.ErrInvalidUserRole
		}
		user.Role = enums.UserRole(role)
	}

	if request.IsActive != nil {
		user.IsActive = *request.IsActive
	}

	if err := s.userRepository.Update(user); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, apperrors.ErrUserNotFound
		}
		return nil, err
	}

	return mapUserToResponse(user), nil
}

func (s *UserService) DeleteUser(id uint) error {
	err := s.userRepository.SoftDelete(id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	return nil
}

func mapUserToResponse(user *domain.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Email:              user.Email,
		Role:               string(user.Role),
		MustChangePassword: user.MustChangePassword,
		IsActive:           user.IsActive,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}
}

func generateTemporaryPassword() (string, error) {
	buffer := make([]byte, 18)
	if _, err := rand.Read(buffer); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer), nil
}
