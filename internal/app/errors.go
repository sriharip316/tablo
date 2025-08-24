package app

import (
	"fmt"
)

// ErrorCode represents different types of application errors
type ErrorCode int

const (
	ErrCodeSuccess ErrorCode = iota
	ErrCodeUsage
	ErrCodeInput
	ErrCodeParse
	ErrCodeProcessing
	ErrCodeSelection
	ErrCodeRender
	ErrCodeOutput
)

// String returns the string representation of the error code
func (e ErrorCode) String() string {
	switch e {
	case ErrCodeSuccess:
		return "success"
	case ErrCodeUsage:
		return "usage_error"
	case ErrCodeInput:
		return "input_error"
	case ErrCodeParse:
		return "parse_error"
	case ErrCodeProcessing:
		return "processing_error"
	case ErrCodeSelection:
		return "selection_error"
	case ErrCodeRender:
		return "render_error"
	case ErrCodeOutput:
		return "output_error"
	default:
		return "unknown_error"
	}
}

// ExitCode returns the appropriate exit code for CLI usage
func (e ErrorCode) ExitCode() int {
	switch e {
	case ErrCodeSuccess:
		return 0
	case ErrCodeUsage:
		return 2
	case ErrCodeInput:
		return 3
	case ErrCodeParse:
		return 4
	case ErrCodeProcessing, ErrCodeSelection, ErrCodeRender, ErrCodeOutput:
		return 5
	default:
		return 1
	}
}

// AppError represents an application-specific error with context
type AppError struct {
	Code    ErrorCode
	Message string
	Cause   error
}

// Error implements the error interface
func (e *AppError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Cause)
	}
	return e.Message
}

// Unwrap returns the underlying cause error
func (e *AppError) Unwrap() error {
	return e.Cause
}

// Is checks if the error matches a specific error code
func (e *AppError) Is(target error) bool {
	if t, ok := target.(*AppError); ok {
		return e.Code == t.Code
	}
	return false
}

// NewError creates a new application error
func NewError(code ErrorCode, message string, cause error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// NewUsageError creates a new usage error
func NewUsageError(message string) *AppError {
	return NewError(ErrCodeUsage, message, nil)
}

// NewInputError creates a new input error
func NewInputError(message string, cause error) *AppError {
	return NewError(ErrCodeInput, message, cause)
}

// NewParseError creates a new parse error
func NewParseError(message string, cause error) *AppError {
	return NewError(ErrCodeParse, message, cause)
}

// NewSelectionError creates a new selection error
func NewSelectionError(message string) *AppError {
	return NewError(ErrCodeSelection, message, nil)
}

// IsUsageError checks if an error is a usage error
func IsUsageError(err error) bool {
	var appErr *AppError
	return AsAppError(err, &appErr) && appErr.Code == ErrCodeUsage
}

// IsInputError checks if an error is an input error
func IsInputError(err error) bool {
	var appErr *AppError
	return AsAppError(err, &appErr) && appErr.Code == ErrCodeInput
}

// IsParseError checks if an error is a parse error
func IsParseError(err error) bool {
	var appErr *AppError
	return AsAppError(err, &appErr) && appErr.Code == ErrCodeParse
}

// AsAppError extracts an AppError from an error chain
func AsAppError(err error, target **AppError) bool {
	for err != nil {
		if appErr, ok := err.(*AppError); ok {
			*target = appErr
			return true
		}
		if unwrapper, ok := err.(interface{ Unwrap() error }); ok {
			err = unwrapper.Unwrap()
		} else {
			break
		}
	}
	return false
}

// GetExitCode extracts the appropriate exit code from an error
func GetExitCode(err error) int {
	if err == nil {
		return 0
	}

	var appErr *AppError
	if AsAppError(err, &appErr) {
		return appErr.Code.ExitCode()
	}

	return 1 // Generic error exit code
}
