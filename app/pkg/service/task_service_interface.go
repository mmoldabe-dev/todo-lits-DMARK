package service

import "todo-lits-DMARK/app/pkg/models"

type TaskService interface {
	CreateTask(req *models.CreateTaskRequest) (*models.Task, error)
	GetTask(id int) (*models.Task, error)
	GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error)
	UpdateTask(id int, updates *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(id int) error
	ToggleTaskStatus(id int) (*models.Task, error)
	GetOverdueTasks() ([]*models.Task, error)
	GetTasksByDateFilter(dateFilter string) ([]*models.Task, error)
	GetTaskStats() (*TaskStats, error)
}
type TaskStats struct {
	Total     int `json:"total"`
	Pending   int `json:"pending"`
	Completed int `json:"completed"`
	Overdue   int `json:"overdue"`
}
