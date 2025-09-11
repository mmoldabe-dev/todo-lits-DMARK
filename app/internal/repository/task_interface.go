package repository

import "todo-lits-DMARK/app/internal/models"

type TaskRepository interface {
	Create(req *models.CreateTaskRequest) (*models.Task, error)
	GetAll(filter *models.TaskFilter, sort *models.SortOption) ([]*models.Task, error)
	GetByID(id int) (*models.Task, error)
	Update(id int, req *models.UpdateTaskRequest) (*models.Task, error)
	Delete(id int) error
	Close() error
}
