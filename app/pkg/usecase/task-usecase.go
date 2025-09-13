package usecase

import (
	"fmt"
	"strings"
	"time"
	"todo-lits-DMARK/app/pkg/models"
	"todo-lits-DMARK/app/pkg/service"
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

func (uc *taskUsecase) GetDashboardData() (*DashboardData, error) {

	stats, err := uc.taskService.GetTaskStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get task stats: %w", err)
	}

	recentSort := &models.TaskSort{Field: "created_at", Order: "desc"}
	allTasks, err := uc.taskService.GetAllTasks(nil, recentSort)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent tasks: %w", err)
	}

	var recentTasks []*models.Task
	if len(allTasks) > 5 {
		recentTasks = allTasks[:5]
	} else {
		recentTasks = allTasks
	}

	overdueTasks, err := uc.taskService.GetOverdueTasks()
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}

	todayTasks, err := uc.taskService.GetTasksByDateFilter("today")
	if err != nil {
		return nil, fmt.Errorf("failed to get today tasks: %w", err)
	}

	upcomingTasks, err := uc.taskService.GetTasksByDateFilter("week")
	if err != nil {
		return nil, fmt.Errorf("failed to get upcoming tasks: %w", err)
	}

	var filteredUpcoming []*models.Task
	today := time.Now().Format("2006-01-02")

	for _, task := range upcomingTasks {
		if task.DueDate != nil && task.DueDate.Format("2006-01-02") != today {
			filteredUpcoming = append(filteredUpcoming, task)
		}
	}

	return &DashboardData{
		Stats:         stats,
		RecentTasks:   recentTasks,
		OverdueTasks:  overdueTasks,
		TodayTasks:    todayTasks,
		UpcomingTasks: filteredUpcoming,
	}, nil
}

func (uc *taskUsecase) SearchTasks(query string) ([]*models.Task, error) {
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	allTasks, err := uc.taskService.GetAllTasks(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to search tasks: %w", err)
	}

	var matchedTasks []*models.Task
	query = strings.ToLower(strings.TrimSpace(query))

	for _, task := range allTasks {

		if strings.Contains(strings.ToLower(task.Title), query) ||
			strings.Contains(strings.ToLower(task.Description), query) {
			matchedTasks = append(matchedTasks, task)
		}
	}

	return matchedTasks, nil
}

func (uc *taskUsecase) BulkUpdateTasks(ids []int, updates *models.UpdateTaskRequest) error {
	if len(ids) == 0 {
		return fmt.Errorf("no task IDs provided")
	}

	if updates.Title != nil {
		title := strings.TrimSpace(*updates.Title)
		if title == "" {
			return fmt.Errorf("task title cannot be empty")
		}
		updates.Title = &title
	}

	if updates.Description != nil {
		description := strings.TrimSpace(*updates.Description)
		updates.Description = &description
	}

	if updates.DueDate != nil && updates.DueDate.Before(time.Now()) {
		return fmt.Errorf("due date cannot be in the past")
	}

	var errors []string
	for _, id := range ids {
		if _, err := uc.taskService.UpdateTask(id, updates); err != nil {
			errors = append(errors, fmt.Sprintf("failed to update task %d: %v", id, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("bulk update failed: %s", strings.Join(errors, "; "))
	}

	return nil
}
