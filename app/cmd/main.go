package main

import (
	"context"
	"log"
	"time"
	"todo-lits-DMARK/app/internal/config"
	"todo-lits-DMARK/app/internal/database"
	"todo-lits-DMARK/app/internal/models"
	"todo-lits-DMARK/app/internal/repository"
	"todo-lits-DMARK/app/internal/service"
	"todo-lits-DMARK/app/internal/usecase"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
)

// App структура для Wails приложения
type App struct {
	ctx         context.Context
	taskUsecase usecase.TaskUsecase
	db          *database.Database
}

// NewApp создает новый экземпляр приложения
func NewApp() *App {
	return &App{}
}

// startup вызывается при запуске приложения
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	// Инициализируем конфигурацию
	cfg := config.New()

	// Подключаемся к базе данных
	db, err := database.New(cfg)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		log.Println("Application will continue without database persistence")
		return
	}

	a.db = db

	// Инициализируем слои
	taskRepo := repository.NewTaskRepository(db.DB)
	taskService := service.NewTaskService(taskRepo)
	a.taskUsecase = usecase.NewTaskUsecase(taskService)

	log.Println("Application started successfully")
}

// shutdown вызывается при завершении приложения
func (a *App) shutdown(ctx context.Context) {
	if a.db != nil {
		if err := a.db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}
}

// CreateTask создает новую задачу
func (a *App) CreateTask(title, description, priority string, dueDate string) (map[string]interface{}, error) {
	req := &models.CreateTaskRequest{
		Title:       title,
		Description: description,
		Priority:    models.TaskPriority(priority),
	}

	// Парсим дату, если она предоставлена
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

// GetTasks получает все задачи
func (a *App) GetTasks(status, priority, sortBy, sortOrder string) ([]map[string]interface{}, error) {
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

// GetTask получает задачу по ID
func (a *App) GetTask(id int) (map[string]interface{}, error) {
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

// UpdateTask обновляет задачу
func (a *App) UpdateTask(id int, title, description, status, priority string, dueDate string) (map[string]interface{}, error) {
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

// DeleteTask удаляет задачу
func (a *App) DeleteTask(id int) error {
	return a.taskUsecase.DeleteTask(id)
}

// ToggleTaskComplete переключает статус задачи
func (a *App) ToggleTaskComplete(id int) (map[string]interface{}, error) {
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

// GetDashboardData получает данные для дашборда
func (a *App) GetDashboardData() (map[string]interface{}, error) {
	data, err := a.taskUsecase.GetDashboardData()
	if err != nil {
		return nil, err
	}

	// Конвертируем задачи в map для JSON
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

// SearchTasks ищет задачи
func (a *App) SearchTasks(query string) ([]map[string]interface{}, error) {
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

// GetTasksByDateFilter получает задачи по фильтру даты
func (a *App) GetTasksByDateFilter(filter string) ([]map[string]interface{}, error) {
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

func main() {
	// Создаем экземпляр приложения
	app := NewApp()

	// Создаем приложение с опциями
	err := wails.Run(&options.App{
		Title:  "TodoApp - Управление задачами",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		OnShutdown:       app.shutdown,
		ContextMenus: &options.ContextMenu{
			Enable: true,
		},
		EnableDefaultContextMenu: false,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
