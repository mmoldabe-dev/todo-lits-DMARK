package repository

import (
	"time"
	"todo-lits-DMARK/app/pkg/models"
)

type TaskRepositoryInterface interface {
	Create(task *models.Task) error
	GetByID(id int) (*models.Task, error)
	GetAll(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error)
	Update(id int, updates *models.UpdateTaskRequest) error
	Delete(id int) error
	GetOverdue() ([]*models.Task, error)
	GetByDateRange(from, to time.Time) ([]*models.Task, error)
}
