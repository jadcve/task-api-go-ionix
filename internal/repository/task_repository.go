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

func (r *TaskRepository) FindAllByAssignedTo(userID uint) ([]domain.Task, error) {
	query := `
		SELECT id, title, description, due_date, status, assigned_to, created_by, created_at, updated_at, deleted_at
		FROM tasks
		WHERE assigned_to = $1 AND deleted_at IS NULL
		ORDER BY id ASC
	`

	rows, err := r.db.Query(context.Background(), query, userID)
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
		    status = $4,
		    assigned_to = $5,
		    updated_at = NOW()
		WHERE id = $6 AND deleted_at IS NULL
		RETURNING updated_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		task.Title,
		task.Description,
		task.DueDate,
		task.Status,
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

func (r *TaskRepository) AddComment(comment *domain.TaskComment) error {
	query := `
		INSERT INTO task_comments (task_id, user_id, comment)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`

	err := r.db.QueryRow(
		context.Background(),
		query,
		comment.TaskID,
		comment.UserID,
		comment.Comment,
	).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *TaskRepository) FindCommentsByTaskID(taskID uint) ([]domain.TaskComment, error) {
	query := `
		SELECT tc.id, tc.task_id, tc.user_id, tc.comment, tc.created_at
		FROM task_comments tc
		INNER JOIN tasks t ON t.id = tc.task_id
		WHERE tc.task_id = $1 AND t.deleted_at IS NULL
		ORDER BY tc.id ASC
	`

	rows, err := r.db.Query(context.Background(), query, taskID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := make([]domain.TaskComment, 0)
	for rows.Next() {
		var comment domain.TaskComment
		if err := rows.Scan(
			&comment.ID,
			&comment.TaskID,
			&comment.UserID,
			&comment.Comment,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}
