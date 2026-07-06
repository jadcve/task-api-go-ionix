package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"task-api-go-ionix/internal/domain"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, must_change_password, is_active, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	var user domain.User
	err := r.db.QueryRow(context.Background(), query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.MustChangePassword,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, must_change_password, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		user.Name,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.MustChangePassword,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, must_change_password, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.MustChangePassword,
		&user.IsActive,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindAll() ([]domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, must_change_password, is_active, created_at, updated_at
		FROM users
		ORDER BY id ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]domain.User, 0)
	for rows.Next() {
		var user domain.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.MustChangePassword,
			&user.IsActive,
			&user.CreatedAt,
			&user.UpdatedAt,
		); err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users
		SET name = $1,
		    role = $2,
		    is_active = $3,
		    updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at
	`

	err := r.db.QueryRow(context.Background(), query, user.Name, user.Role, user.IsActive, user.ID).Scan(&user.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *UserRepository) UpdatePassword(userID uint, passwordHash string, mustChangePassword bool) error {
	query := `
		UPDATE users
		SET password_hash = $1,
		    must_change_password = $2,
		    updated_at = NOW()
		WHERE id = $3
	`

	result, err := r.db.Exec(context.Background(), query, passwordHash, mustChangePassword, userID)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *UserRepository) SoftDelete(id uint) error {
	query := `
		UPDATE users
		SET is_active = FALSE,
		    updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.Exec(context.Background(), query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}
