package handlers

import (
	"myCalculator/internal/domain"
	"myCalculator/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type CalculationHandler struct {
	service service.CalculationService
}

func NewCalculationHandler(s service.CalculationService) *CalculationHandler {
	return &CalculationHandler{service: s}
}

// Основные методы ORM - Create, Read(Find), Update, Delete.
func (h *CalculationHandler) GetCalculation(c echo.Context) error {
	calculations, err := h.service.GetAllCalculations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not get calculations"})
	}

	return c.JSON(http.StatusOK, calculations)
}

func (h *CalculationHandler) PostCalculation(c echo.Context) error {
	var req domain.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	calc, err := h.service.CreateCalculation(req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not create calculation"})
	}

	return c.JSON(http.StatusCreated, calc)
}

func (h *CalculationHandler) PatchCalculation(c echo.Context) error {
	id := c.Param("id")

	var req domain.CalculationRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	updatedCalc, err := h.service.UpdateCalculation(id, req.Expression)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Could not update calculation"})
	}

	return c.JSON(http.StatusOK, updatedCalc)
}

func (h *CalculationHandler) DeleteCalculation(c echo.Context) error {
	id := c.Param("id")

	if err := h.service.DeleteCalculation(id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Could not delete calculations"})
	}

	return c.NoContent(http.StatusNoContent)
}
