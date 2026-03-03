package main

import (
	"fmt"
	"net/http"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Calculation struct {
	ID         string `json:"id"`
	Expression string `json:"expression"`
	Result     string `json:"result"`
}

type CalculationRequest struct {
	Expression string `json:"expression"`
}

type TaskRequest struct {
	Task string `json:"task"`
}

type Task struct {
	ID     string `json:"id"`
	Task   string `json:"task"`
	Status string `json:"status"`
}

type TaskUpdateRequest struct {
	Task   string `json:"task"`
	Status string `json:"status"`
}

var calculations = []Calculation{}

var tasks = []Task{}
var task string

func getHelloHandler(c echo.Context) error {
	if task == "" {
		return c.JSON(http.StatusOK, "hello,")
	}

	return c.JSON(http.StatusOK, "hello"+" "+task)
}

func postTaskHandler(c echo.Context) error {
	var req TaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	if req.Task == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Task cannot be empty",
		})
	}

	newTask := Task{
		ID:     uuid.NewString(),
		Task:   req.Task,
		Status: "pending",
	}

	tasks = append(tasks, newTask)
	task = req.Task

	return c.JSON(http.StatusCreated, newTask)
}

func getTasksHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, tasks)
}

// НОВЫЙ: GET задачи по ID
func getTaskByIDHandler(c echo.Context) error {
	id := c.Param("id")

	for _, task := range tasks {
		if task.ID == id {
			return c.JSON(http.StatusOK, task)
		}
	}

	return c.JSON(http.StatusNotFound, map[string]string{
		"error": "Task not found",
	})
}

// НОВЫЙ: PATCH handler для обновления task по ID
func patchTaskHandler(c echo.Context) error {
	id := c.Param("id")

	// Ищем задачу по ID
	var taskToUpdate *Task
	for i := range tasks {
		if tasks[i].ID == id {
			taskToUpdate = &tasks[i]
			break
		}
	}

	if taskToUpdate == nil {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	// Парсим запрос
	var req TaskUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request format",
		})
	}

	// Обновляем поля (если они переданы)
	if req.Task != "" {
		taskToUpdate.Task = req.Task
		task = req.Task // обновляем и старую переменную
	}

	if req.Status != "" {
		taskToUpdate.Status = req.Status
	}

	return c.JSON(http.StatusOK, taskToUpdate)
}

// НОВЫЙ: DELETE handler для удаления task по ID
func deleteTaskHandler(c echo.Context) error {
	id := c.Param("id")

	// Ищем индекс задачи
	index := -1
	for i, task := range tasks {
		if task.ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "Task not found",
		})
	}

	// Удаляем задачу из среза
	deletedTask := tasks[index]
	tasks = append(tasks[:index], tasks[index+1:]...)

	// Если удалили текущую задачу, обновляем currentTask
	if deletedTask.Task == task {
		if len(tasks) > 0 {
			task = tasks[len(tasks)-1].Task // берем последнюю
		} else {
			task = ""
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task deleted successfully",
		"id":      id,
	})
}

// func patchTaskHandler(c echo.Context) error {
// 	id := c.Param("id")

// 	var req TaskRequest
// 	if err := c.Bind(&req); err != nil {
// 		return c.JSON(http.StatusBadRequest, map[string]string{
// 			"error": "Invalid request format",
// 		})
// 	}

// 	if req.Task == id {
// 		task = req.Task
// 		return c.JSON(http.StatusOK, map[string]string{
// 			"message": "Task updated successfully",
// 			"task":    task,
// 		})
// 	}
// 	return c.JSON(http.StatusBadRequest, map[string]string{
// 		"error": "Task cannot be empty",
// 	})
// }

// func deleteTaskHandler(c echo.Context) error {
// 	id := c.Param("id")

// 	if task == id {
// 		return c.NoContent(http.StatusNoContent)
// 	}
// 	return c.JSON(http.StatusNoContent, task)
// }

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

func getCalculation(c echo.Context) error {
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

	calculations = append(calculations, calc)

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

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations[i].Expression = req.Expression
			calculations[i].Result = result
		}
	}
	return c.JSON(http.StatusOK, calculations)
}

func deleteCalculation(c echo.Context) error {
	id := c.Param("id")

	for i, calculation := range calculations {
		if calculation.ID == id {
			calculations = append(calculations[:i], calculations[i+1:]...)
		}
	}
	return c.NoContent(http.StatusNoContent)
}

func main() {

	e := echo.New()

	e.Use(middleware.CORS())
	e.Use(middleware.RequestLogger())

	e.GET("/calculations", getCalculation)
	e.POST("/calculations", postCalculation)
	e.PATCH("/calculations/:id", patchCalculation)
	e.DELETE("/calculations/:id", deleteCalculation)

	e.POST("/task", postTaskHandler)
	e.GET("/", getHelloHandler)
	// НОВЫЕ эндпоинты для работы с задачами по ID
	e.GET("/tasks", getTasksHandler)          // получить все задачи
	e.GET("/tasks/:id", getTaskByIDHandler)   // получить задачу по ID
	e.PATCH("/tasks/:id", patchTaskHandler)   // ОБНОВИТЬ задачу по ID
	e.DELETE("/tasks/:id", deleteTaskHandler) // УДАЛИТЬ задачу по ID

	e.Start("localhost:8080")
}
