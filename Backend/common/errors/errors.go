package errors

import "fmt"

type AppError struct {
	Code    string
	Message string
}

func (e AppError) Error() string {
	return e.Message
}

func New(code, message string) AppError {
	return AppError{Code: code, Message: message}
}

func Wrap(code string, err error) AppError {
	return AppError{Code: code, Message: fmt.Sprintf("%v", err)}
}
