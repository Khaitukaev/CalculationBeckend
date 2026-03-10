package handlers

import (
	"myCalculator/internal/domain"
	"myCalculator/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type TaskHandler struct {
	service service.TaskService
}

func NewTaskHandler(service service.TaskService) *TaskHandler {
	return &TaskHandler{service: service}
}

// CreateTask godoc
// @Summary Create a new task
// @Tags tasks
// @Accept json
// @Produce json
// @Success 201 {object} domain.Task
func (h *TaskHandler) CreateTask(c echo.Context) error {
	var req domain.CreateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	task, err := h.service.CreateTask(&req)
	if err != nil {
		switch err.(type) {
		case domain.ValidationError:
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not create task"})
		}
	}

	return c.JSON(http.StatusCreated, task)
}

// GetAllTasks godoc
// @Summary Get all tasks
// @Tags tasks
// @Produce json
// @Param includeDeleted query bool false "Include soft deleted tasks"
// @Success 200 {array} domain.Task
func (h *TaskHandler) GetAllTasks(c echo.Context) error {
	includeDeleted := c.QueryParam("includeDeleted") == "true"

	tasks, err := h.service.GetAllTasks(includeDeleted)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch tasks"})
	}

	return c.JSON(http.StatusOK, tasks)
}

// GetTaskByID godoc
// @Summary Get task by ID
// @Tags tasks
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
func (h *TaskHandler) GetTaskByID(c echo.Context) error {
	id := c.Param("id")

	task, err := h.service.GetTaskByID(id)
	if err != nil {
		switch err.(type) {
		case domain.NotFoundError:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
	}

	return c.JSON(http.StatusOK, task)
}

// UpdateTask godoc
// @Summary Update task
// @Tags tasks
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} domain.Task
func (h *TaskHandler) UpdateTask(c echo.Context) error {
	id := c.Param("id")

	var req domain.UpdateTaskRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	task, err := h.service.UpdateTask(id, &req)
	if err != nil {
		switch err.(type) {
		case domain.NotFoundError:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update task"})
		}
	}

	return c.JSON(http.StatusOK, task)
}

// DeleteTask godoc
// @Summary Delete task
// @Tags tasks
// @Param id path string true "Task ID"
// @Success 200 {object} map[string]string
func (h *TaskHandler) DeleteTask(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteTask(id); err != nil {
		switch err.(type) {
		case domain.NotFoundError:
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		default:
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete task"})
		}
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task soft deleted successfully",
		"id":      id,
	})
}

// GetHello godoc
// @Summary Get hello with current task
// @Tags tasks
// @Produce plain
// @Success 200 {string} string
func (h *TaskHandler) GetHello(c echo.Context) error {
	currentTask, err := h.service.GetCurrentTask()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get current task"})
	}

	if currentTask == "" {
		return c.String(http.StatusOK, "hello")
	}
	return c.String(http.StatusOK, "hello "+currentTask)
}
