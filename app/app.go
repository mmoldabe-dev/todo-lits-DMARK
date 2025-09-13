package app

import (
	"context"
	"log"
	"time"
	"todo-lits-DMARK/app/pkg/config"
	"todo-lits-DMARK/app/pkg/database"
	"todo-lits-DMARK/app/pkg/models"
	"todo-lits-DMARK/app/pkg/repository"
	"todo-lits-DMARK/app/pkg/service"
	"todo-lits-DMARK/app/pkg/usecase"
)

type App struct {
	ctx         context.Context
	taskUsecase usecase.TaskUsecase
	db          *database.Database
}

func NewApp() *App {
	return &App{}
}

func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx
	log.Println("TodoApp is starting...")

	cfg := config.New()
	log.Printf("Connecting to database: %s:%s", cfg.Database.Host, cfg.Database.Port)

	db, err := database.New(cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Application will continue without database persistence")
		return
	}

	a.db = db
	log.Println("Database connection established")

	taskRepo := repository.NewTaskRepository(db.DB)
	taskService := service.NewTaskService(taskRepo)
	a.taskUsecase = usecase.NewTaskUsecase(taskService)

	log.Println("Application started successfully")
}

func (a *App) OnShutdown(ctx context.Context) {
	log.Println("TodoApp is shutting down...")
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		} else {
			log.Println("Database connection closed")
		}
	}
}

// Greet для тестирования Wails
func (a *App) Greet(name string) string {
	return "Hello " + name + " from TodoApp!"
}

func (a *App) CreateTask(title, description, priority string, dueDate string) (map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return nil, nil // Graceful fallback если DB недоступна
	}

	req := &models.CreateTaskRequest{
		Title:       title,
		Description: description,
		Priority:    models.TaskPriority(priority),
	}

	if dueDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05Z", dueDate); err == nil {
			req.DueDate = &parsedDate
		}
	}

	task, err := a.taskUsecase.CreateTask(req)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
		"is_overdue":  task.IsOverdue(),
	}, nil
}

func (a *App) GetTasks(status, priority, sortBy, sortOrder string) ([]map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return []map[string]interface{}{}, nil
	}

	tasks, err := a.taskUsecase.GetTasks(status, priority, sortBy, sortOrder)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		result[i] = map[string]interface{}{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
			"is_overdue":  task.IsOverdue(),
		}
	}

	return result, nil
}

func (a *App) GetTask(id int) (map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return nil, nil
	}

	task, err := a.taskUsecase.GetTask(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
		"is_overdue":  task.IsOverdue(),
	}, nil
}

func (a *App) UpdateTask(id int, title, description, status, priority string, dueDate string) (map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return nil, nil
	}

	updates := &models.UpdateTaskRequest{}

	if title != "" {
		updates.Title = &title
	}
	if description != "" {
		updates.Description = &description
	}
	if status != "" {
		taskStatus := models.TaskStatus(status)
		updates.Status = &taskStatus
	}
	if priority != "" {
		taskPriority := models.TaskPriority(priority)
		updates.Priority = &taskPriority
	}
	if dueDate != "" {
		if parsedDate, err := time.Parse("2006-01-02T15:04:05Z", dueDate); err == nil {
			updates.DueDate = &parsedDate
		}
	}

	task, err := a.taskUsecase.UpdateTask(id, updates)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
		"is_overdue":  task.IsOverdue(),
	}, nil
}

func (a *App) DeleteTask(id int) error {
	if a.taskUsecase == nil {
		return nil
	}
	return a.taskUsecase.DeleteTask(id)
}

func (a *App) ToggleTaskComplete(id int) (map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return nil, nil
	}

	task, err := a.taskUsecase.ToggleTaskComplete(id)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"id":          task.ID,
		"title":       task.Title,
		"description": task.Description,
		"status":      task.Status,
		"priority":    task.Priority,
		"due_date":    task.DueDate,
		"created_at":  task.CreatedAt,
		"updated_at":  task.UpdatedAt,
		"is_overdue":  task.IsOverdue(),
	}, nil
}

func (a *App) GetDashboardData() (map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return map[string]interface{}{
			"stats": map[string]interface{}{
				"total": 0, "pending": 0, "completed": 0, "overdue": 0,
			},
			"recent_tasks":   []map[string]interface{}{},
			"overdue_tasks":  []map[string]interface{}{},
			"today_tasks":    []map[string]interface{}{},
			"upcoming_tasks": []map[string]interface{}{},
		}, nil
	}

	data, err := a.taskUsecase.GetDashboardData()
	if err != nil {
		return nil, err
	}

	convertTasks := func(tasks []*models.Task) []map[string]interface{} {
		result := make([]map[string]interface{}, len(tasks))
		for i, task := range tasks {
			result[i] = map[string]interface{}{
				"id":          task.ID,
				"title":       task.Title,
				"description": task.Description,
				"status":      task.Status,
				"priority":    task.Priority,
				"due_date":    task.DueDate,
				"created_at":  task.CreatedAt,
				"updated_at":  task.UpdatedAt,
				"is_overdue":  task.IsOverdue(),
			}
		}
		return result
	}

	return map[string]interface{}{
		"stats": map[string]interface{}{
			"total":     data.Stats.Total,
			"pending":   data.Stats.Pending,
			"completed": data.Stats.Completed,
			"overdue":   data.Stats.Overdue,
		},
		"recent_tasks":   convertTasks(data.RecentTasks),
		"overdue_tasks":  convertTasks(data.OverdueTasks),
		"today_tasks":    convertTasks(data.TodayTasks),
		"upcoming_tasks": convertTasks(data.UpcomingTasks),
	}, nil
}

func (a *App) SearchTasks(query string) ([]map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return []map[string]interface{}{}, nil
	}

	tasks, err := a.taskUsecase.SearchTasks(query)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		result[i] = map[string]interface{}{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
			"is_overdue":  task.IsOverdue(),
		}
	}

	return result, nil
}

func (a *App) GetTasksByDateFilter(filter string) ([]map[string]interface{}, error) {
	if a.taskUsecase == nil {
		return []map[string]interface{}{}, nil
	}

	tasks, err := a.taskUsecase.GetTasksByDateRange(filter)
	if err != nil {
		return nil, err
	}

	result := make([]map[string]interface{}, len(tasks))
	for i, task := range tasks {
		result[i] = map[string]interface{}{
			"id":          task.ID,
			"title":       task.Title,
			"description": task.Description,
			"status":      task.Status,
			"priority":    task.Priority,
			"due_date":    task.DueDate,
			"created_at":  task.CreatedAt,
			"updated_at":  task.UpdatedAt,
			"is_overdue":  task.IsOverdue(),
		}
	}

	return result, nil
}
