package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	dsn := "host=localhost user=postgres password=mypasswoed dbname=postgres port=5432 sslmode=disable"
	var err error

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Could not connect to database: %v", err)
	}

	if err := db.AutoMigrate(&Calculation{}, Task{}); err != nil {
		log.Fatalf("Could not migrate: %v", err)
	}

}

type Calculation struct {
	ID         string `gorm:"primaryKey "json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

// Task модель для базы данных с мягким удалением
type Task struct {
	ID        string         `gorm:primaryKey json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`                           // вместо Status
	CreatedAt time.Time      `json:"created_at"`                        // автоматически заполняется GORM
	UpdatedAt time.Time      `json:"updated_at"`                        // автоматически обновляется GORM
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"` // для мягкого удаления
}

// TaskRequest для создания задачи (с полями task и is_done)
type TaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

// TaskUpdateRequest для обновления задачи
type TaskUpdateRequest struct {
	Task   *string `json:"task,omitempty"` // указатель для опциональных полей
	IsDone *bool   `json:"is_done,omitempty"`
}

var currentTask string

func getHelloHandler(c echo.Context) error {
	if currentTask == "" {
		// Если нет текущей задачи, берем последнюю из БД
		var lastTask Task
		db.Last(&lastTask)
		if lastTask.ID != "" {
			currentTask = lastTask.Task
		} else {
			return c.JSON(http.StatusOK, "hello")
		}
	}
	return c.JSON(http.StatusOK, "hello "+currentTask)
}

// 1. POST /tasks - создание задачи в БД
func postTaskHandler(c echo.Context) error {
	var req TaskRequest

	// Парсим JSON
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Проверяем, что задача не пустая
	if req.Task == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Task cannot be empty",
		})
	}

	// Создаем задачу для БД
	newTask := Task{
		ID:     uuid.NewString(),
		Task:   req.Task,
		IsDone: req.IsDone,
	}

	// Сохраняем в БД
	if err := db.Create(&newTask).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add tasks"})
	}

	// Обновляем currentTask для обратной совместимости
	currentTask = req.Task

	return c.JSON(http.StatusCreated, newTask)
}

// 2. GET /tasks - получение всех задач из БД (включая мягко удаленные)
func getTasksHandler(c echo.Context) error {
	var tasks []Task

	// Unscoped() показывает даже мягко удаленные записи
	if err := db.Unscoped().Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tasks"})
	}

	return c.JSON(http.StatusOK, tasks)
}

// GET /tasks/active - только активные (не удаленные) задачи
func getActiveTasksHandler(c echo.Context) error {
	var tasks []Task

	// Без Unscoped() показываются только записи с NULL в deleted_at
	if err := db.Find(&tasks).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to fetch tasks",
		})
	}

	return c.JSON(http.StatusOK, tasks)
}

// 3. GET /tasks/:id - получение задачи по ID
func getTaskByIDHandler(c echo.Context) error {
	id := c.Param("id")

	var task Task

	// Ищем в БД (включая мягко удаленные, чтобы можно было их восстановить)
	if err := db.Unscoped().First(&task, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Task not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	return c.JSON(http.StatusOK, task)
}

// 4. PATCH /tasks/:id - обновление задачи в БД
func patchTaskHandler(c echo.Context) error {
	id := c.Param("id")

	// Ищем задачу по ID
	var task Task

	if err := db.First(&task, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Task not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// Парсим запрос
	var req TaskUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Обновляем поля (только те, что переданы)
	updates := make(map[string]interface{})

	if req.Task != nil {
		updates["task"] = *req.Task
		currentTask = *req.Task // обновляем текущую задачу
	}

	if req.IsDone != nil {
		updates["is_done"] = *req.IsDone
	}

	// Если есть что обновлять
	if len(updates) > 0 {
		updates["updated_at"] = time.Now()

		if err := db.Model(&task).Updates(updates).Error; err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"error": "Failed to update task",
			})
		}
	}

	// Получаем обновленную задачу
	db.First(&task, "id = ?", id)

	return c.JSON(http.StatusOK, task)
}

// 5. DELETE /tasks/:id - мягкое удаление задачи из БД
func deleteTaskHandler(c echo.Context) error {
	id := c.Param("id")

	// Ищем задачу
	var task Task
	if err := db.First(&task, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{
				"error": "Task not found",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
	}

	// Мягкое удаление (ставит DeletedAt)
	if err := db.Delete(&task).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to delete task",
		})
	}

	// Если удалили текущую задачу, обновляем currentTask
	if task.Task == currentTask {
		var lastTask Task
		db.Last(&lastTask)
		if lastTask.ID != "" {
			currentTask = lastTask.Task
		} else {
			currentTask = ""
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task soft deleted successfully",
		"id":      id,
	})
}

func calculateExpression(expression string) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(expression)
	if err != nil {
		return "", err
	}

	result, err := expr.Evaluate(nil)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%v", result), err
}

// Основные методы ORM - Create, Read(Find), Update, Delete.
func getCalculation(c echo.Context) error {
	var calculations []Calculation

	if err := db.Find(&calculations).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

func postCalculation(c echo.Context) error {
	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	calc := Calculation{
		ID:         uuid.NewString(),
		Expression: req.Expression,
		Result:     result,
	}

	if err := db.Create(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not add calculations"})
	}

	return c.JSON(http.StatusCreated, calc)
}

func patchCalculation(c echo.Context) error {
	id := c.Param("id")

	var req CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	result, err := calculateExpression(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid expression"})
	}

	var calc Calculation
	if err := db.First(&calc, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not find expression"})
	}

	calc.Expression = req.Expression
	calc.Result = result

	if err := db.Save(&calc).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not update calculations"})
	}

	return c.JSON(http.StatusOK, calc)
}

func deleteCalculation(c echo.Context) error {
	id := c.Param("id")

	if err := db.Delete(&Calculation{}, "id = ?", id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculations"})
	}

	return c.NoContent(http.StatusNoContent)
}

func main() {
	initDB()

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	e.GET("/calculations", getCalculation)
	e.POST("/calculations", postCalculation)
	e.PATCH("/calculations/:id", patchCalculation)
	e.DELETE("/calculations/:id", deleteCalculation)

	// Эндпоинты для задач
	e.POST("/task", postTaskHandler) // старый POST для совместимости
	e.GET("/", getHelloHandler)      // старый GET

	// НОВЫЕ эндпоинты для работы с БД
	e.POST("/tasks", postTaskHandler)             // создать задачу
	e.GET("/tasks", getTasksHandler)              // все задачи (с удаленными)
	e.GET("/tasks/active", getActiveTasksHandler) // только активные
	e.GET("/tasks/:id", getTaskByIDHandler)       // получить по ID
	e.PATCH("/tasks/:id", patchTaskHandler)       // обновить
	e.DELETE("/tasks/:id", deleteTaskHandler)     // мягко удалить

	e.Logger.Fatal(e.Start("localhost:8080"))
}
