package service

import (
	"fmt"
	"time"
	"todo-lits-DMARK/app/pkg/models"
	"todo-lits-DMARK/app/pkg/repository"

	"github.com/go-playground/validator/v10"
)

type taskService struct {
	repo      repository.TaskRepositoryInterface
	validator *validator.Validate
}

func NewTaskService(repo repository.TaskRepositoryInterface) TaskService {
	validator := validator.New()

	validator.RegisterValidation("task_status", validateTaskStatus)
	validator.RegisterValidation("task_priority", validateTaskPriority)

	return &taskService{
		repo:      repo,
		validator: validator,
	}
}
func validateTaskStatus(fl validator.FieldLevel) bool {
	status := fl.Field().String()
	return status == string(models.TaskStatusPending) || status == string(models.TaskStatusCompleted)
}

func validateTaskPriority(fl validator.FieldLevel) bool {
	priority := fl.Field().String()
	return priority == string(models.TaskPriorityLow) ||
		priority == string(models.TaskPriorityMedium) ||
		priority == string(models.TaskPriorityHigh)
}

func (s *taskService) CreateTask(req *models.CreateTaskRequest) (*models.Task, error) {
	if err := s.validator.Struct(req); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if req.Priority == "" {
		req.Priority = models.TaskPriorityMedium
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Priority:    req.Priority,
		DueDate:     req.DueDate,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, fmt.Errorf("failed to create task: %w", err)
	}

	return task, nil
}
func (s *taskService) GetTask(id int) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID: %d", id)
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return task, nil
}
func (s *taskService) GetAllTasks(filter *models.TaskFilter, sort *models.TaskSort) ([]*models.Task, error) {
	if sort != nil {
		if err := s.validator.Struct(sort); err != nil {
			return nil, fmt.Errorf("invalid sort parameters: %w", err)
		}
	}

	tasks, err := s.repo.GetAll(filter, sort)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}

	return tasks, nil
}
func (s *taskService) UpdateTask(id int, updates *models.UpdateTaskRequest) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID: %d", id)
	}

	if err := s.validator.Struct(updates); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if _, err := s.repo.GetByID(id); err != nil {
		return nil, fmt.Errorf("task not found: %w", err)
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, fmt.Errorf("failed to update task: %w", err)
	}

	return s.repo.GetByID(id)
}
func (s *taskService) DeleteTask(id int) error {
	if id <= 0 {
		return fmt.Errorf("invalid task ID: %d", id)
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	return nil
}
func (s *taskService) ToggleTaskStatus(id int) (*models.Task, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid task ID: %d", id)
	}

	task, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	var newStatus models.TaskStatus
	if task.Status == models.TaskStatusPending {
		newStatus = models.TaskStatusCompleted
	} else {
		newStatus = models.TaskStatusPending
	}

	updates := &models.UpdateTaskRequest{
		Status: &newStatus,
	}

	if err := s.repo.Update(id, updates); err != nil {
		return nil, fmt.Errorf("failed to toggle task status: %w", err)
	}

	return s.repo.GetByID(id)
}
func (s *taskService) GetOverdueTasks() ([]*models.Task, error) {
	tasks, err := s.repo.GetOverdue()
	if err != nil {
		return nil, fmt.Errorf("failed to get overdue tasks: %w", err)
	}

	return tasks, nil
}
func (s *taskService) GetTasksByDateFilter(dateFilter string) ([]*models.Task, error) {
	now := time.Now()
	var from, to time.Time

	switch dateFilter {
	case "today":
		from = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		to = from.Add(24 * time.Hour).Add(-time.Second)
	case "week":

		weekday := int(now.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		from = now.AddDate(0, 0, -(weekday - 1))
		from = time.Date(from.Year(), from.Month(), from.Day(), 0, 0, 0, 0, from.Location())
		to = from.AddDate(0, 0, 7).Add(-time.Second)
	case "overdue":
		return s.GetOverdueTasks()
	default:
		return nil, fmt.Errorf("invalid date filter: %s", dateFilter)
	}

	tasks, err := s.repo.GetByDateRange(from, to)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks by date filter: %w", err)
	}

	return tasks, nil
}
func (s *taskService) GetTaskStats() (*TaskStats, error) {
	allTasks, err := s.repo.GetAll(nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks for stats: %w", err)
	}

	stats := &TaskStats{}
	stats.Total = len(allTasks)

	for _, task := range allTasks {
		switch task.Status {
		case models.TaskStatusPending:
			stats.Pending++
			if task.IsOverdue() {
				stats.Overdue++
			}
		case models.TaskStatusCompleted:
			stats.Completed++
		}
	}

	return stats, nil
}
