package repository

import (
	"myCalculator/internal/domain"

	"gorm.io/gorm"
)

// Основные методы CRUD - Create, Read, Update, Delete

type CalculationRepository interface {
	CreateCalculation(calc *domain.Calculation) error
	GetAllCalculations() ([]domain.Calculation, error)
	GetCalculationByID(id string) (domain.Calculation, error)
	UpdateCalculation(calc domain.Calculation) error
	DeleteCalculation(id string) error
}

type calcRepository struct {
	db *gorm.DB
}

func NewCalculationRepository(db *gorm.DB) CalculationRepository {
	return &calcRepository{db: db}
}

func (r *calcRepository) CreateCalculation(calc *domain.Calculation) error {
	return r.db.Create(&calc).Error
}

func (r *calcRepository) GetAllCalculations() ([]domain.Calculation, error) {
	var calculations []domain.Calculation
	err := r.db.Find(&calculations).Error
	return calculations, err
}

func (r *calcRepository) GetCalculationByID(id string) (domain.Calculation, error) {
	var calc domain.Calculation
	err := r.db.First(&calc, "id = ?", id).Error
	return calc, err
}

func (r *calcRepository) UpdateCalculation(calc domain.Calculation) error {
	return r.db.Save(&calc).Error
}

func (r *calcRepository) DeleteCalculation(id string) error {
	return r.db.Delete(&domain.Calculation{}, "id = ?", id).Error
}
