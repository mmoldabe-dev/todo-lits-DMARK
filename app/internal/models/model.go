package models

import "time"

type Task struct {
	ID          int        `json:"id" db:"id"`
	Title       string     `json:"title" db:"title"`
	Description string     `json:"description" db:"description"`
	Priority    Priority   `json:"priority" db:"priority"`
	Completed   bool       `json:"completed" db:"completed"`
	DueDate     *time.Time `json:"due_date" db:"due_date"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
}

func (t *Task) IsOverDue() bool {
	if t.DueDate == nil || t.Completed {
		return false
	}
	return t.DueDate.Before(time.Now())
}

type CreateTaskRequest struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Priority    Priority   `json:"priority"`
	DueDate     *time.Time `json:"due_date"`
}
type UpdateTaskRequest struct {
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Priority    *Priority  `json:"priority"`
	Completed   *bool      `json:"completed"`
	DueDate     *time.Time `json:"due_date"`
}
type TaskFilter struct {
	Status     string `json:"status"`
	Priority   string `json:"priority"`
	DateFilter string `json:"date_filter"`
}
type SortOption struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

func (t *Task)IsDueToday()bool{
	if t.DueDate == nil || t.Completed{
		return  false
	}
	now := time.Now()
	due := *t.DueDate

	return now.Year() == due.Year() && now.YearDay() == due.YearDay()
}

func(t*Task)IsDueWeek()bool{
	if t.DueDate == nil || t.Completed{
		return  false
	}
	now := time.Now()
	due := *t.DueDate

	return now.Year() == due.Year() && now.Month() == due.Month() &&  now.YearDay() == due.YearDay()
}