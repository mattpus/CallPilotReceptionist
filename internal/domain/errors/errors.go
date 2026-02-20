package errors

import "fmt"

type DomainError struct {
	Code    string
	Message string
	Err     error
}

func (e *DomainError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func (e *DomainError) Unwrap() error {
	return e.Err
}

// Error codes
const (
	ErrCodeNotFound        = "NOT_FOUND"
	ErrCodeAlreadyExists   = "ALREADY_EXISTS"
	ErrCodeInvalidInput    = "INVALID_INPUT"
	ErrCodeUnauthorized    = "UNAUTHORIZED"
	ErrCodeForbidden       = "FORBIDDEN"
	ErrCodeInternal        = "INTERNAL_ERROR"
	ErrCodeProviderError   = "PROVIDER_ERROR"
	ErrCodeDatabaseError   = "DATABASE_ERROR"
	ErrCodeValidationError = "VALIDATION_ERROR"
)

// Common domain errors
func NewNotFoundError(entity string, id string) *DomainError {
	return &DomainError{
		Code:    ErrCodeNotFound,
		Message: fmt.Sprintf("%s with id %s not found", entity, id),
	}
}

func NewAlreadyExistsError(entity string, field string, value string) *DomainError {
	return &DomainError{
		Code:    ErrCodeAlreadyExists,
		Message: fmt.Sprintf("%s with %s=%s already exists", entity, field, value),
	}
}

func NewInvalidInputError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeInvalidInput,
		Message: message,
	}
}

func NewValidationError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeValidationError,
		Message: message,
	}
}

func NewUnauthorizedError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeUnauthorized,
		Message: message,
	}
}

func NewForbiddenError(message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeForbidden,
		Message: message,
	}
}

func NewInternalError(err error) *DomainError {
	return &DomainError{
		Code:    ErrCodeInternal,
		Message: "Internal server error",
		Err:     err,
	}
}

func NewProviderError(err error, message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeProviderError,
		Message: message,
		Err:     err,
	}
}

func NewDatabaseError(err error, message string) *DomainError {
	return &DomainError{
		Code:    ErrCodeDatabaseError,
		Message: message,
		Err:     err,
	}
}
