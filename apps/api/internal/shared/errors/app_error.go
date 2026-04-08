package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
	Status  int
	Cause   error
}

func New(code string, message string, status int) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
	}
}

func Wrap(code string, message string, status int, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Cause:   cause,
	}
}

func (appError *AppError) Error() string {
	if appError.Cause == nil {
		return fmt.Sprintf("%s: %s", appError.Code, appError.Message)
	}

	return fmt.Sprintf("%s: %s: %v", appError.Code, appError.Message, appError.Cause)
}
