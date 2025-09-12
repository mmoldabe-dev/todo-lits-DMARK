package usecase

import (
	"fmt"
	"strings"
	"time"
	"todo-lits-DMARK/app/internal/models"
	"todo-lits-DMARK/app/internal/service"
)

type TaskUsecase interface {
	// Основные операции CRUD
	CreateTask(req *models.CreateTaskRequest) (*models.Task, error)
	GetTask(id int) (*models.Task, error)
	GetTasks(status string, priority string, sortBy string, sortOrder string) ([]*models.Task, error)
	UpdateTask(id int, updates *models.UpdateTaskRequest) (*models.Task, error)
	DeleteTask(id int) error

	// Специальные операции
	ToggleTaskComplete(id int) (*models.Task, error)
	GetTasksByDateRange(dateFilter string) ([]*models.Task, error)
	GetDashboardData() (*DashboardData, error)
	SearchTasks(query string) ([]*models.Task, error)
	BulkUpdateTasks(ids []int, updates *models.UpdateTaskRequest) error
}
type DashboardData struct {
	Stats         *service.TaskStats `json:"stats"`
	RecentTasks   []*models.Task     `json:"recent_tasks"`
	OverdueTasks  []*models.Task     `json:"overdue_tasks"`
	TodayTasks    []*models.Task     `json:"today_tasks"`
	UpcomingTasks []*models.Task     `json:"upcoming_tasks"`
}
type taskUsecase struct {
	taskService service.TaskService
}

func NewTaskUsecase(taskService service.TaskService) TaskUsecase {
	return &taskUsecase{
		taskService: taskService,
	}
}
func (uc *taskUsecase) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {

	if strings.TrimSpace(req.Title) == "" {
		return nil, fmt.Errorf("task title cannot be empty")
	}

	req.Title = strings.TrimSpace(req.Title)
	req.Description = strings.TrimSpace(req.Description)

	if req.DueDate != nil && req.DueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	task, err := uc.taskService.CreateTask(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
func (uc *taskUsecase) GetTask(id int) (*models.Task, error) {
	return uc.taskService.GetTask(id)
}


func (uc *taskUsecase) GetTasks(status string, priority string, sortBy string, sortOrder string) ([]*models.Task, error) {

	filter := &models.TaskFilter{}

	if status != "" && status != "all" {
		taskStatus := models.TaskStatus(status)
		filter.Status = &taskStatus
	}

	if priority != "" && priority != "all" {
		taskPriority := models.TaskPriority(priority)
		filter.Priority = &taskPriority
	}

	var sort *models.TaskSort
	if sortBy != "" {
		sort = &models.TaskSort{
			Field: sortBy,
			Order: "asc",
		}

		if sortOrder == "desc" {
			sort.Order = "desc"
		}
	}

	return uc.taskService.GetAllTasks(filter, sort)
}


func (uc *taskUsecase) UpdateTask(id int, updates *models.UpdateTaskRequest) (*models.Task, error) {
	
	if updates.Title != nil {
		title := strings.TrimSpace(*updates.Title)
		if title == "" {
			return nil, fmt.Errorf("task title cannot be empty")
		}
		updates.Title = &title
	}

	if updates.Description != nil {
		description := strings.TrimSpace(*updates.Description)
		updates.Description = &description
	}

	if updates.DueDate != nil && updates.DueDate.Before(time.Now()) {
		return nil, fmt.Errorf("due date cannot be in the past")
	}

	return uc.taskService.UpdateTask(id, updates)
}


func (uc *taskUsecase) DeleteTask(id int) error {
	return uc.taskService.DeleteTask(id)
}


func (uc *taskUsecase) ToggleTaskComplete(id int) (*models.Task, error) {
	return uc.taskService.ToggleTaskStatus(id)
}


func (uc *taskUsecase) GetTasksByDateRange(dateFilter string) ([]*models.Task, error) {
	validFilters := map[string]bool{
		"today":   true,
		"week":    true,
		"overdue": true,
	}

	if !validFilters[dateFilter] {
		return nil, fmt.Errorf("invalid date filter: %s. Valid filters: today, week, overdue", dateFilter)
	}

	return uc.taskService.GetTasksByDateFilter(dateFilter)
}


