package service

import (
	"myCalculator/internal/domain"
	"myCalculator/internal/repository"
	"time"
)

type TaskService interface {
	CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error)
	GetAllTasks(includeDeleted bool) ([]domain.Task, error)
	GetTaskByID(id string) (*domain.Task, error)
	UpdateTask(id string, req *domain.UpdateTaskRequest) (*domain.Task, error)
	DeleteTask(id string) error
	GetCurrentTask() (string, error)
}

type taskService struct {
	repo        repository.TaskRepository
	currentTask string
}

func NewTaskService(repo repository.TaskRepository) TaskService {
	return &taskService{repo: repo}
}

func (s *taskService) CreateTask(req *domain.CreateTaskRequest) (*domain.Task, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	task := domain.NewTask(req.Task, req.IsDone)

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	s.currentTask = task.Task
	return task, nil
}

func (s *taskService) GetAllTasks(includeDeleted bool) ([]domain.Task, error) {
	return s.repo.GetAll(includeDeleted)
}

func (s *taskService) GetTaskByID(id string) (*domain.Task, error) {
	return s.repo.GetByID(id)
}

func (s *taskService) UpdateTask(id string, req *domain.UpdateTaskRequest) (*domain.Task, error) {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	updates := make(map[string]interface{})

	if req.Task != nil {
		updates["task"] = *req.Task
		s.currentTask = *req.Task
	}

	if req.IsDone != nil {
		updates["is_done"] = *req.IsDone
	}

	if len(updates) > 0 {
		updates["updated_at"] = time.Now()
		if err := s.repo.Update(id, updates); err != nil {
			return nil, err
		}
	}

	// Получаем обновленную задачу
	return s.repo.GetByID(id)
}

func (s *taskService) DeleteTask(id string) error {
	task, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return err
	}

	// Обновляем currentTask если удалили текущую
	if task.Task == s.currentTask {
		lastTask, _ := s.repo.GetLast()
		if lastTask != nil {
			s.currentTask = lastTask.Task
		} else {
			s.currentTask = ""
		}
	}

	return nil
}

func (s *taskService) GetCurrentTask() (string, error) {
	if s.currentTask != "" {
		return s.currentTask, nil
	}

	lastTask, err := s.repo.GetLast()
	if err != nil {
		return "", err
	}

	if lastTask != nil {
		s.currentTask = lastTask.Task
		return lastTask.Task, nil
	}

	return "", nil
}
