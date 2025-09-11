package repository

import "todo-lits-DMARK/app/internal/models"

type TaskRepositoryInterface interface {
	Create(task *models.Task) error
	GetAll() ([]models.Task, error)
	GetByID(id int) (*models.Task, error)
	Update(task *models.Task) error
	Delete(id int) error
	GetByStatus(completed bool) ([]models.Task, error)
}
