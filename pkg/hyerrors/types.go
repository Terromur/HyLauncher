package hyerrors

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

// ErrorType represents the category of error
type ErrorType string

const (
	ErrorTypeGame       ErrorType = "GAME"
	ErrorTypeJava       ErrorType = "JAVA"
	ErrorTypeNetwork    ErrorType = "NETWORK"
	ErrorTypeValidation ErrorType = "VALIDATION"
	ErrorTypeFileSystem ErrorType = "FILESYSTEM"
	ErrorTypeConfig     ErrorType = "CONFIG"
	ErrorTypeUpdate     ErrorType = "UPDATE"
	ErrorTypeUnknown    ErrorType = "UNKNOWN"
)

// AppError represents an application error with context
type AppError struct {
	Type      ErrorType `json:"type"`
	Message   string    `json:"message"`
	Technical string    `json:"technical"`
	Timestamp time.Time `json:"timestamp"`
	Stack     string    `json:"stack"`
}

// ErrorHandler is an interface for handling errors
type ErrorHandler interface {
	HandleError(err *AppError)
}

var globalHandler ErrorHandler

// SetHandler sets the global error handler
func SetHandler(handler ErrorHandler) {
	globalHandler = handler
}

// Error implements the error interface
func (e *AppError) Error() string {
	return e.Message
}

// NewAppError creates a new application error
func NewAppError(errType ErrorType, userMessage string, err error) *AppError {
	technical := ""
	if err != nil {
		technical = err.Error()
	}

	// Capture stack trace
	stack := captureStack(3)

	appErr := &AppError{
		Type:      errType,
		Message:   userMessage,
		Technical: technical,
		Timestamp: time.Now(),
		Stack:     stack,
	}

	// Notify handler if set
	if globalHandler != nil {
		globalHandler.HandleError(appErr)
	}

	return appErr
}

// captureStack captures the call stack
func captureStack(skip int) string {
	stack := ""
	for i := skip; i < skip+10; i++ {
		pc, file, line, ok := runtime.Caller(i)
		if !ok {
			break
		}
		fn := runtime.FuncForPC(pc)
		stack += fmt.Sprintf("%s:%d %s\n", filepath.Base(file), line, fn.Name())
	}
	return stack
}

// IsCritical returns true if the error is critical and should trigger crash reporting
func (e *AppError) IsCritical() bool {
	return e.Type == ErrorTypeGame || e.Type == ErrorTypeFileSystem
}
