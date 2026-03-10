package repository

import (
	"myCalculator/internal/domain"

	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	GetAll(includeDeleted bool) ([]domain.Task, error)
	GetByID(id string) (*domain.Task, error)
	Update(id string, updates map[string]interface{}) error
	Delete(id string) error
	GetLast() (*domain.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) GetAll(includeDeleted bool) ([]domain.Task, error) {
	var tasks []domain.Task
	query := r.db
	if includeDeleted {
		query = query.Unscoped()
	}
	err := query.Find(&tasks).Error
	return tasks, err
}

func (r *taskRepository) GetByID(id string) (*domain.Task, error) {
	var task domain.Task
	err := r.db.Unscoped().First(&task, "id = ?", id).Error
	if err == gorm.ErrRecordNotFound {
		return nil, domain.NewNotFoundError("task not found")
	}
	return &task, err
}

func (r *taskRepository) Update(id string, updates map[string]interface{}) error {
	return r.db.Model(&domain.Task{}).Where("id = ?", id).Updates(updates).Error
}

func (r *taskRepository) Delete(id string) error {
	return r.db.Delete(&domain.Task{}, "id = ?", id).Error
}

func (r *taskRepository) GetLast() (*domain.Task, error) {
	var task domain.Task
	err := r.db.Last(&task).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return &task, err
}
