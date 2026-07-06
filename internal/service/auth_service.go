package service

import (
	"errors"

	"github.com/jackc/pgx/v5"

	"task-api-go-ionix/internal/domain"
	"task-api-go-ionix/internal/dto"
	apperrors "task-api-go-ionix/internal/errors"
	"task-api-go-ionix/internal/security"
)

type UserRepository interface {
	FindByEmail(email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	UpdatePassword(userID uint, passwordHash string, mustChangePassword bool) error
}

type AuthService struct {
	userRepository     UserRepository
	jwtSecret          string
	jwtExpirationHours int
}

func NewAuthService(userRepository UserRepository, jwtSecret string, jwtExpirationHours int) *AuthService {
	return &AuthService{
		userRepository:     userRepository,
		jwtSecret:          jwtSecret,
		jwtExpirationHours: jwtExpirationHours,
	}
}

func (s *AuthService) Login(request dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrInvalidCredentials
	}
	if !user.IsActive {
		return nil, apperrors.ErrInactiveUser
	}
	if !security.CheckPassword(request.Password, user.PasswordHash) {
		return nil, apperrors.ErrInvalidCredentials
	}

	token, err := security.GenerateToken(user.ID, string(user.Role), s.jwtSecret, s.jwtExpirationHours)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.AuthUserResponse{
			ID:                 user.ID,
			Name:               user.Name,
			Email:              user.Email,
			Role:               string(user.Role),
			MustChangePassword: user.MustChangePassword,
		},
	}, nil
}

func (s *AuthService) ChangePassword(userID uint, request dto.ChangePasswordRequest) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return apperrors.ErrUserNotFound
	}
	if !security.CheckPassword(request.CurrentPassword, user.PasswordHash) {
		return apperrors.ErrInvalidCurrentPass
	}

	hash, err := security.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	err = s.userRepository.UpdatePassword(userID, hash, false)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return apperrors.ErrUserNotFound
		}
		return err
	}

	return nil
}

func (s *AuthService) Me(userID uint) (*dto.MeResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, apperrors.ErrUserNotFound
	}
	if !user.IsActive {
		return nil, apperrors.ErrInactiveUser
	}

	return &dto.MeResponse{
		ID:                 user.ID,
		Name:               user.Name,
		Email:              user.Email,
		Role:               string(user.Role),
		MustChangePassword: user.MustChangePassword,
	}, nil
}
