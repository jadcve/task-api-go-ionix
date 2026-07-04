package repository

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"task-api-go-ionix/internal/domain"
)

type TaskRepository struct {
	db *pgxpool.Pool
}

func NewTaskRepository(db *pgxpool.Pool) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) Create(task *domain.Task) error {
	query := `
		INSERT INTO tasks (title, description, due_date, status, assigned_to, created_by)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Status,
		task.AssignedTo,
		task.CreatedBy,
	).Scan(&task.ID, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) FindByID(id uint) (*domain.Task, error) {
	query := `
		SELECT id, title, description, due_date, status, assigned_to, created_by, created_at, updated_at, deleted_at
		FROM tasks
		WHERE id = $1 AND deleted_at IS NULL
	`

	var task domain.Task
	var deletedAt *time.Time
	err := r.db.QueryRow(context.Background(), query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.DueDate,
		&task.Status,
		&task.AssignedTo,
		&task.CreatedBy,
		&task.CreatedAt,
		&task.UpdatedAt,
		&deletedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	task.DeletedAt = deletedAt
	return &task, nil
}

func (r *TaskRepository) FindAll() ([]domain.Task, error) {
	query := `
		SELECT id, title, description, due_date, status, assigned_to, created_by, created_at, updated_at, deleted_at
		FROM tasks
		WHERE deleted_at IS NULL
		ORDER BY id ASC
	`

	rows, err := r.db.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tasks := make([]domain.Task, 0)
	for rows.Next() {
		var task domain.Task
		var deletedAt *time.Time
		if err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.DueDate,
			&task.Status,
			&task.AssignedTo,
			&task.CreatedBy,
			&task.CreatedAt,
			&task.UpdatedAt,
			&deletedAt,
		); err != nil {
			return nil, err
		}
		task.DeletedAt = deletedAt
		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *TaskRepository) Update(task *domain.Task) error {
	query := `
		UPDATE tasks
		SET title = $1,
		    description = $2,
		    due_date = $3,
		    assigned_to = $4,
		    updated_at = NOW()
		WHERE id = $5 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.AssignedTo,
		task.ID,
	).Scan(&task.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return pgx.ErrNoRows
		}
		return err
	}

	return nil
}

func (r *TaskRepository) Delete(id uint) error {
	query := `
		UPDATE tasks
		SET deleted_at = NOW(),
		    updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
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

func (r *TaskRepository) FindAssignedUser(userID uint) (*domain.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, must_change_password, is_active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	var user domain.User
	err := r.db.QueryRow(context.Background(), query, userID).Scan(
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
