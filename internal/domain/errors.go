package domain

type ValidationError struct {
	Message string
}

func (e ValidationError) Error() string {
	return e.Message
}

func NewValidationError(message string) ValidationError {
	return ValidationError{Message: message}
}

type NotFoundError struct {
	Message string
}

func (e NotFoundError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) NotFoundError {
	return NotFoundError{Message: message}
}
