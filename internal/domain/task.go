package domain

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Task struct {
	ID        string         `gorm:"primaryKey" json:"id"`
	Task      string         `json:"task"`
	IsDone    bool           `json:"is_done"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

func NewTask(task string, isDone bool) *Task {
	return &Task{
		ID:     uuid.New().String(),
		Task:   task,
		IsDone: isDone,
	}
}

func (t *Task) Update(task *string, isDone *bool) {
	if task != nil {
		t.Task = *task
	}
	if isDone != nil {
		t.IsDone = *isDone
	}
	t.UpdatedAt = time.Now()
}
