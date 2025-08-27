package main

import (
	"testing"

	"github.com/sriharip316/tablo/internal/app"
)

func TestOptionsToAppConfig(t *testing.T) {
	opts := &options{
		file:   "test.json",
		inStr:  `{"a": 1}`,
		format: "json",
		dive:   true,
		style:  "ascii",
		quiet:  true,
	}

	config := opts.toAppConfig()

	if config.Input.File != "test.json" {
		t.Errorf("expected file to be test.json, got %s", config.Input.File)
	}
	if config.Input.String != `{"a": 1}` {
		t.Errorf("expected string to be {\"a\": 1}, got %s", config.Input.String)
	}
	if config.Input.Format != "json" {
		t.Errorf("expected format to be json, got %s", config.Input.Format)
	}
	if !config.Flatten.Enabled {
		t.Error("expected flatten to be enabled")
	}
	if config.Output.Style != "ascii" {
		t.Errorf("expected style to be ascii, got %s", config.Output.Style)
	}
	if !config.General.Quiet {
		t.Error("expected quiet to be true")
	}
}

func TestHandleError(t *testing.T) {
	// Test that error handling works with different error types
	usageErr := app.NewUsageError("test usage error")
	exitCode := app.GetExitCode(usageErr)
	if exitCode != 2 {
		t.Errorf("expected exit code 2 for usage error, got %d", exitCode)
	}

	inputErr := app.NewInputError("test input error", nil)
	exitCode = app.GetExitCode(inputErr)
	if exitCode != 3 {
		t.Errorf("expected exit code 3 for input error, got %d", exitCode)
	}
}
