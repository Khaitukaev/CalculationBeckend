package domain

type CreateTaskRequest struct {
	Task   string `json:"task"`
	IsDone bool   `json:"is_done"`
}

func (r *CreateTaskRequest) Validate() error {
	if r.Task == "" {
		return ErrEmptyTask
	}
	return nil
}

type UpdateTaskRequest struct {
	Task   *string `json:"task,omitempty"`
	IsDone *bool   `json:"is_done,omitempty"`
}

var ErrEmptyTask = NewValidationError("task cannot be empty")
