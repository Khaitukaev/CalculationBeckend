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

var calculations = []Calculation{}
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

	task = req.Task

	// Убедитесь, что этот return есть и он правильный
	return c.JSON(http.StatusOK, map[string]string{
		"message": "Task updated successfully",
		"task":    task,
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

	e.POST("/task", postTaskHandler)
	e.GET("/", getHelloHandler)

	e.Start("localhost:8080")
}
