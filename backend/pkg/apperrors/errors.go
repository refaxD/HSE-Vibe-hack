package apperrors

import "fmt"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func BadRequest(msg string) *AppError {
	return &AppError{Code: 400, Message: msg}
}

func Unauthorized() *AppError {
	return &AppError{Code: 401, Message: "Unauthorized"}
}

func NotFound(entity string) *AppError {
	return &AppError{Code: 404, Message: fmt.Sprintf("%s not found", entity)}
}

func Conflict(msg string) *AppError {
	return &AppError{Code: 409, Message: msg}
}

func UserAlreadyExists() *AppError {
	return &AppError{Code: 409, Message: "User already exists"}
}

func UserNotFound() *AppError {
	return &AppError{Code: 409, Message: "User not found"}
}

func InvalidPassword(msg string) *AppError {
	return &AppError{Code: 409, Message: fmt.Sprintf("Invalid User Password: %s", msg)}
}

func ValidationError(errors map[string]string) *AppError {
	return &AppError{Code: 409, Message: fmt.Sprintf("%v", errors)}
}

func InternalServer(msg ...string) *AppError {
	m := "Internal Server Error"
	if len(msg) > 0 {
		m = msg[0]
	}
	return &AppError{Code: 500, Message: m}
}
