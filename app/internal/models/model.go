package models

import "time"

type TaskStatus string

const (
	TaskStatusPending   TaskStatus = "pending"
	TaskStatusCompleted TaskStatus = "completed"
)

type TaskPriority string

const (
	TaskPriorityLow    TaskPriority = "low"
	TaskPriorityMedium TaskPriority = "medium"
	TaskPriorityHigh   TaskPriority = "high"
)

type Task struct {
	ID          int          `json:"id" db:"id"`
	Title       string       `json:"title" db:"title" validate:"required,min=1,max=255"`
	Description string       `json:"description" db:"description"`
	Status      TaskStatus   `json:"status" db:"status"`
	Priority    TaskPriority `json:"priority" db:"priority"`
	DueDate     *time.Time   `json:"due_date" db:"due_date"`
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
}
type CreateTaskRequest struct {
	Title       string       `json:"title" validate:"required,min=1,max=255"`
	Description string       `json:"description" validate:"max=1000"`
	Priority    TaskPriority `json:"priority" validate:"oneof=low medium high"`
	DueDate     *time.Time   `json:"due_date"`
}
type UpdateTaskRequest struct {
	Title       *string       `json:"title,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string       `json:"description,omitempty" validate:"omitempty,max=1000"`
	Status      *TaskStatus   `json:"status,omitempty" validate:"omitempty,oneof=pending completed"`
	Priority    *TaskPriority `json:"priority,omitempty" validate:"omitempty,oneof=low medium high"`
	DueDate     *time.Time    `json:"due_date,omitempty"`
}
type TaskFilter struct {
	Status   *TaskStatus   `json:"status,omitempty"`
	Priority *TaskPriority `json:"priority,omitempty"`
	DateFrom *time.Time    `json:"date_from,omitempty"`
	DateTo   *time.Time    `json:"date_to,omitempty"`
}
type TaskSort struct {
	Field string `json:"field" validate:"oneof=created_at due_date priority"`
	Order string `json:"order" validate:"oneof=asc desc"`
}

func (t *Task) IsOverdue() bool {
	if t.DueDate == nil || t.Status == TaskStatusCompleted {
		return false
	}
	return t.DueDate.Before(time.Now())
}
func (t *Task) PriorityValue() int {
	switch t.Priority {
	case TaskPriorityHigh:
		return 3
	case TaskPriorityMedium:
		return 2
	case TaskPriorityLow:
		return 1
	default:
		return 0
	}
}
