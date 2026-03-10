package service

import (
	"fmt"
	"myCalculator/internal/domain"
	"myCalculator/internal/repository"

	"github.com/Knetic/govaluate"
	"github.com/google/uuid"
)

type CalculationService interface {
	CreateCalculation(expression string) (domain.Calculation, error)
	GetAllCalculations() ([]domain.Calculation, error)
	GetCalculationByID(id string) (domain.Calculation, error)
	UpdateCalculation(id, expression string) (domain.Calculation, error)
	DeleteCalculation(id string) error
}

type calcService struct {
	repo repository.CalculationRepository
}

func NewCalculationService(r repository.CalculationRepository) CalculationService {
	return &calcService{repo: r}
}

func (s *calcService) calculateExpression(expression string) (string, error) {
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

// CreateCalculation implements [CalculationService].
func (s *calcService) CreateCalculation(expression string) (domain.Calculation, error) {
	result, err := s.calculateExpression(expression)
	if err != nil {
		return domain.Calculation{}, err
	}

	calc := domain.Calculation{
		ID:         uuid.NewString(),
		Expression: expression,
		Result:     result,
	}

	if err := s.repo.CreateCalculation(calc.Result); err != nil {
		return domain.Calculation{}, err
	}

	return calc, nil
}

// GetAllCalculations implements [CalculationService].
func (s *calcService) GetAllCalculations() ([]domain.Calculation, error) {
	return s.repo.GetAllCalculations()
}

// GetCalculationByID implements [CalculationService].
func (s *calcService) GetCalculationByID(id string) (domain.Calculation, error) {
	return s.repo.GetCalculationByID(id)
}

// UpdateCalculation implements [CalculationService].
func (s *calcService) UpdateCalculation(id string, expression string) (domain.Calculation, error) {
	calc, err := s.repo.GetCalculationByID(id)
	if err != nil {
		return domain.Calculation{}, err
	}

	result, err := s.calculateExpression(expression)
	if err != nil {
		return domain.Calculation{}, err
	}

	calc.Expression = expression
	calc.Result = result

	if err := s.repo.UpdateCalculation(calc); err != nil {
		return domain.Calculation{}, err
	}

	return calc, nil
}

// DeleteCalculation implements [CalculationService].
func (s *calcService) DeleteCalculation(id string) error {
	return s.repo.DeleteCalculation(id)
}
