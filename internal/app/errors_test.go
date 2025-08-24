package app

import (
	"errors"
	"testing"
)

func TestErrorCode_String(t *testing.T) {
	testCases := []struct {
		code     ErrorCode
		expected string
	}{
		{ErrCodeSuccess, "success"},
		{ErrCodeUsage, "usage_error"},
		{ErrCodeInput, "input_error"},
		{ErrCodeParse, "parse_error"},
		{ErrCodeProcessing, "processing_error"},
		{ErrCodeSelection, "selection_error"},
		{ErrCodeRender, "render_error"},
		{ErrCodeOutput, "output_error"},
		{ErrorCode(999), "unknown_error"},
	}

	for _, tc := range testCases {
		result := tc.code.String()
		if result != tc.expected {
			t.Errorf("ErrorCode(%d).String(): got %s, want %s", tc.code, result, tc.expected)
		}
	}
}

func TestErrorCode_ExitCode(t *testing.T) {
	testCases := []struct {
		code     ErrorCode
		expected int
	}{
		{ErrCodeSuccess, 0},
		{ErrCodeUsage, 2},
		{ErrCodeInput, 3},
		{ErrCodeParse, 4},
		{ErrCodeProcessing, 5},
		{ErrCodeSelection, 5},
		{ErrCodeRender, 5},
		{ErrCodeOutput, 5},
		{ErrorCode(999), 1},
	}

	for _, tc := range testCases {
		result := tc.code.ExitCode()
		if result != tc.expected {
			t.Errorf("ErrorCode(%d).ExitCode(): got %d, want %d", tc.code, result, tc.expected)
		}
	}
}

func TestAppError_Error(t *testing.T) {
	// Test error without cause
	err := &AppError{
		Code:    ErrCodeUsage,
		Message: "invalid flag",
		Cause:   nil,
	}
	expected := "invalid flag"
	if err.Error() != expected {
		t.Errorf("AppError.Error() without cause: got %s, want %s", err.Error(), expected)
	}

	// Test error with cause
	causeErr := errors.New("underlying error")
	err = &AppError{
		Code:    ErrCodeInput,
		Message: "failed to read",
		Cause:   causeErr,
	}
	expected = "failed to read: underlying error"
	if err.Error() != expected {
		t.Errorf("AppError.Error() with cause: got %s, want %s", err.Error(), expected)
	}
}

func TestAppError_Unwrap(t *testing.T) {
	causeErr := errors.New("underlying error")
	err := &AppError{
		Code:    ErrCodeInput,
		Message: "failed to read",
		Cause:   causeErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != causeErr {
		t.Errorf("AppError.Unwrap(): got %v, want %v", unwrapped, causeErr)
	}

	// Test without cause
	err = &AppError{
		Code:    ErrCodeUsage,
		Message: "invalid flag",
		Cause:   nil,
	}

	unwrapped = err.Unwrap()
	if unwrapped != nil {
		t.Errorf("AppError.Unwrap() without cause: got %v, want nil", unwrapped)
	}
}

func TestAppError_Is(t *testing.T) {
	err1 := &AppError{Code: ErrCodeUsage, Message: "test"}
	err2 := &AppError{Code: ErrCodeUsage, Message: "different message"}
	err3 := &AppError{Code: ErrCodeInput, Message: "test"}
	regularErr := errors.New("regular error")

	if !err1.Is(err2) {
		t.Error("AppError.Is(): expected true for same error code")
	}

	if err1.Is(err3) {
		t.Error("AppError.Is(): expected false for different error code")
	}

	if err1.Is(regularErr) {
		t.Error("AppError.Is(): expected false for non-AppError")
	}
}

func TestNewError(t *testing.T) {
	causeErr := errors.New("cause")
	err := NewError(ErrCodeInput, "test message", causeErr)

	if err.Code != ErrCodeInput {
		t.Errorf("NewError code: got %v, want %v", err.Code, ErrCodeInput)
	}
	if err.Message != "test message" {
		t.Errorf("NewError message: got %s, want test message", err.Message)
	}
	if err.Cause != causeErr {
		t.Errorf("NewError cause: got %v, want %v", err.Cause, causeErr)
	}
}

func TestNewUsageError(t *testing.T) {
	err := NewUsageError("invalid usage")

	if err.Code != ErrCodeUsage {
		t.Errorf("NewUsageError code: got %v, want %v", err.Code, ErrCodeUsage)
	}
	if err.Message != "invalid usage" {
		t.Errorf("NewUsageError message: got %s, want invalid usage", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("NewUsageError cause: got %v, want nil", err.Cause)
	}
}

func TestNewInputError(t *testing.T) {
	causeErr := errors.New("file not found")
	err := NewInputError("cannot read file", causeErr)

	if err.Code != ErrCodeInput {
		t.Errorf("NewInputError code: got %v, want %v", err.Code, ErrCodeInput)
	}
	if err.Message != "cannot read file" {
		t.Errorf("NewInputError message: got %s, want cannot read file", err.Message)
	}
	if err.Cause != causeErr {
		t.Errorf("NewInputError cause: got %v, want %v", err.Cause, causeErr)
	}
}

func TestNewParseError(t *testing.T) {
	causeErr := errors.New("invalid JSON")
	err := NewParseError("parse failed", causeErr)

	if err.Code != ErrCodeParse {
		t.Errorf("NewParseError code: got %v, want %v", err.Code, ErrCodeParse)
	}
	if err.Message != "parse failed" {
		t.Errorf("NewParseError message: got %s, want parse failed", err.Message)
	}
	if err.Cause != causeErr {
		t.Errorf("NewParseError cause: got %v, want %v", err.Cause, causeErr)
	}
}

func TestNewSelectionError(t *testing.T) {
	err := NewSelectionError("missing paths")

	if err.Code != ErrCodeSelection {
		t.Errorf("NewSelectionError code: got %v, want %v", err.Code, ErrCodeSelection)
	}
	if err.Message != "missing paths" {
		t.Errorf("NewSelectionError message: got %s, want missing paths", err.Message)
	}
	if err.Cause != nil {
		t.Errorf("NewSelectionError cause: got %v, want nil", err.Cause)
	}
}

func TestIsUsageError(t *testing.T) {
	usageErr := NewUsageError("test")
	inputErr := NewInputError("test", nil)
	regularErr := errors.New("test")

	if !IsUsageError(usageErr) {
		t.Error("IsUsageError: expected true for usage error")
	}

	if IsUsageError(inputErr) {
		t.Error("IsUsageError: expected false for input error")
	}

	if IsUsageError(regularErr) {
		t.Error("IsUsageError: expected false for regular error")
	}
}

func TestIsInputError(t *testing.T) {
	inputErr := NewInputError("test", nil)
	usageErr := NewUsageError("test")
	regularErr := errors.New("test")

	if !IsInputError(inputErr) {
		t.Error("IsInputError: expected true for input error")
	}

	if IsInputError(usageErr) {
		t.Error("IsInputError: expected false for usage error")
	}

	if IsInputError(regularErr) {
		t.Error("IsInputError: expected false for regular error")
	}
}

func TestIsParseError(t *testing.T) {
	parseErr := NewParseError("test", nil)
	usageErr := NewUsageError("test")
	regularErr := errors.New("test")

	if !IsParseError(parseErr) {
		t.Error("IsParseError: expected true for parse error")
	}

	if IsParseError(usageErr) {
		t.Error("IsParseError: expected false for usage error")
	}

	if IsParseError(regularErr) {
		t.Error("IsParseError: expected false for regular error")
	}
}

func TestAsAppError(t *testing.T) {
	appErr := NewInputError("test", nil)
	regularErr := errors.New("regular")

	// Test with AppError
	var target *AppError
	if !AsAppError(appErr, &target) {
		t.Error("AsAppError: expected true for AppError")
	}
	if target != appErr {
		t.Error("AsAppError: target should be set to the AppError")
	}

	// Test with regular error
	target = nil
	if AsAppError(regularErr, &target) {
		t.Error("AsAppError: expected false for regular error")
	}
	if target != nil {
		t.Error("AsAppError: target should remain nil for regular error")
	}

	// Test with nil error
	target = nil
	if AsAppError(nil, &target) {
		t.Error("AsAppError: expected false for nil error")
	}
	if target != nil {
		t.Error("AsAppError: target should remain nil for nil error")
	}
}

func TestAsAppError_WithWrapping(t *testing.T) {
	// Create a wrapped error
	appErr := NewInputError("inner", nil)
	wrappedErr := &wrapperError{err: appErr}

	var target *AppError
	if !AsAppError(wrappedErr, &target) {
		t.Error("AsAppError: expected true for wrapped AppError")
	}
	if target != appErr {
		t.Error("AsAppError: target should be set to the wrapped AppError")
	}
}

func TestGetExitCode(t *testing.T) {
	testCases := []struct {
		err      error
		expected int
	}{
		{nil, 0},
		{NewUsageError("test"), 2},
		{NewInputError("test", nil), 3},
		{NewParseError("test", nil), 4},
		{NewSelectionError("test"), 5},
		{errors.New("regular error"), 1},
	}

	for _, tc := range testCases {
		result := GetExitCode(tc.err)
		if result != tc.expected {
			t.Errorf("GetExitCode(%v): got %d, want %d", tc.err, result, tc.expected)
		}
	}
}

// Helper type for testing error unwrapping
type wrapperError struct {
	err error
}

func (w *wrapperError) Error() string {
	return "wrapped: " + w.err.Error()
}

func (w *wrapperError) Unwrap() error {
	return w.err
}
