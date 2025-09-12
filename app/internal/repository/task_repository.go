package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
	"todo-lits-DMARK/app/internal/models"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) TaskRepositoryInterface {
	return &TaskRepository{
		db: db,
	}
}

func (r *TaskRepository) Create(task *models.Task) error {
	query := `
        INSERT INTO tasks (title, description, status, priority, due_date, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
        RETURNING id
    `

	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = models.TaskStatusPending

	err := r.db.QueryRow(
		query,
		task.Title,
		task.Description,
		task.Status,
		task.Priority,
		task.DueDate,
		task.CreatedAt,
		task.UpdatedAt,
	).Scan(&task.ID)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	return nil
}
func (r *TaskRepository) GetByID(id int) (*models.Task, error) {
	task := &models.Task{}

	query := `
       SELECT id, title, description, status, priority, due_date, created_at, updated_at
        FROM tasks
        WHERE id = $1
    `

	row := r.db.QueryRow(query, id)
	err := row.Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.Priority,
		&task.DueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with id %d not found", id)
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}
func (r *TaskRepository) GetAll(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error) {
	query := `	SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks`

	var tasks []*models.Task
	var conditions []string
	var args []interface{}
	argCount := 0
	//filtesr
	if filter != nil {
		if filter.Status != nil {
			argCount++
			conditions = append(conditions, fmt.Sprintf("status = $%d", argCount))
			args = append(args, *filter.Status)
		}

		if filter.Priority != nil {
			argCount++
			conditions = append(conditions, fmt.Sprintf("priority = $%d", argCount))
			args = append(args, *filter.Priority)
		}
		if filter.DateFrom != nil {
			argCount++
			conditions = append(conditions, fmt.Sprintf("due_date >= $%d", argCount))
			args = append(args, *filter.DateFrom)
		}

		if filter.DateTo != nil {
			argCount++
			conditions = append(conditions, fmt.Sprintf("due_date <= $%d", argCount))
			args = append(args, *filter.DateTo)
		}
	}
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if sort != nil && sort.Field != "" {
		orderBy := sort.Field
		if sort.Order == "desc" {
			orderBy += " DESC"
		} else {
			orderBy += " ASC"
		}

		if sort.Field == "priority" {

			orderBy = "CASE priority WHEN 'high' THEN 3 WHEN 'medium' THEN 2 WHEN 'low' THEN 1 END"
			if sort.Order == "desc" {
				orderBy += " DESC"
			} else {
				orderBy += " ASC"
			}
		}

		query += " ORDER BY " + orderBy
	} else {
		query += " ORDER BY created_at DESC"
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return tasks, nil
}
func (r *TaskRepository) Update(id int, updates *models.UpdateTaskRequest) error {
	var setParts []string
	var args []interface{}
	argCount := 0

	if updates.Title != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("title = $%d", argCount))
		args = append(args, *updates.Title)
	}

	if updates.Description != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("description = $%d", argCount))
		args = append(args, *updates.Description)
	}

	if updates.Status != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("status = $%d", argCount))
		args = append(args, *updates.Status)
	}

	if updates.Priority != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("priority = $%d", argCount))
		args = append(args, *updates.Priority)
	}

	if updates.DueDate != nil {
		argCount++
		setParts = append(setParts, fmt.Sprintf("due_date = $%d", argCount))
		args = append(args, *updates.DueDate)
	}

	if len(setParts) == 0 {
		return fmt.Errorf("no fields to update")
	}

	argCount++
	setParts = append(setParts, fmt.Sprintf("updated_at = $%d", argCount))
	args = append(args, time.Now())

	argCount++
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE tasks SET %s WHERE id = $%d",
		strings.Join(setParts, ", "),
		argCount,
	)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
func (r *TaskRepository) Delete(id int) error {
	query := "DELETE FROM tasks WHERE id = $1"

	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}
func (r *TaskRepository) GetOverdue() ([]*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE due_date < NOW() AND status = 'pending'
		ORDER BY due_date ASC`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}
func (r *TaskRepository) GetByDateRange(from, to time.Time) ([]*models.Task, error) {
	query := `
		SELECT id, title, description, status, priority, due_date, created_at, updated_at
		FROM tasks
		WHERE due_date BETWEEN $1 AND $2
		ORDER BY due_date ASC`

	rows, err := r.db.Query(query, from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by date range: %w", err)
	}
	defer rows.Close()

	var tasks []*models.Task

	for rows.Next() {
		task := &models.Task{}
		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.DueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan task: %w", err)
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return tasks, nil
}
