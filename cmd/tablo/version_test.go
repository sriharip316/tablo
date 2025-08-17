package main

import (
	"errors"
	"os"
	"os/exec"
	"sync"
	"testing"
)

// Helper to reset the global version state for each test.
func resetVersionState() {
	version = ""
	versionOnce = sync.Once{}
}

// Test resolveVersion logic with direct mocking of execCommand.
func TestResolveVersion(t *testing.T) {
	tests := []struct {
		name     string
		setup    func() // Function to set up the execCommand mock
		expected string
	}{
		{
			name: "Injected version takes precedence",
			setup: func() {
				version = "v1.0.0-injected" // This will be used first
				// Ensure that no git commands are actually run if injected version exists
				execCommand = func(name string, arg ...string) *exec.Cmd {
					t.Fatalf("git command %s %v should not be called when version is injected", name, arg)
					return nil // Unreachable
				}
			},
			expected: "v1.0.0-injected",
		},
		{
			name: "Latest tag - clean working tree",
			setup: func() {
				execCommand = func(name string, arg ...string) *exec.Cmd {
					switch {
					case name == "git" && len(arg) > 0 && arg[0] == "describe":
						cmd := exec.Command("echo", "v0.3.0")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "status":
						cmd := exec.Command("echo", "") // Clean working tree
						cmd.Env = os.Environ()
						return cmd
					}
					return exec.Command(name, arg...) // Fallback for other commands if any
				}
			},
			expected: "v0.3.0",
		},
		{
			name: "Latest tag - dirty working tree",
			setup: func() {
				execCommand = func(name string, arg ...string) *exec.Cmd {
					switch {
					case name == "git" && len(arg) > 0 && arg[0] == "describe":
						cmd := exec.Command("echo", "v0.3.0")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "status":
						cmd := exec.Command("echo", " M somefile.txt") // Dirty working tree
						cmd.Env = os.Environ()
						return cmd
					}
					return exec.Command(name, arg...)
				}
			},
			expected: "v0.3.0-dirty",
		},
		{
			name: "Short hash - clean working tree (no tag)",
			setup: func() {
				execCommand = func(name string, arg ...string) *exec.Cmd {
					switch {
					case name == "git" && len(arg) > 0 && arg[0] == "describe":
						// Simulate git describe failure (no tag)
						cmd := exec.Command("false")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "rev-parse":
						cmd := exec.Command("echo", "abcdefg")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "status":
						cmd := exec.Command("echo", "") // Clean
						cmd.Env = os.Environ()
						return cmd
					}
					return exec.Command(name, arg...)
				}
			},
			expected: "dev-abcdefg",
		},
		{
			name: "Short hash - dirty working tree (no tag)",
			setup: func() {
				execCommand = func(name string, arg ...string) *exec.Cmd {
					switch {
					case name == "git" && len(arg) > 0 && arg[0] == "describe":
						// Simulate git describe failure (no tag)
						cmd := exec.Command("false")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "rev-parse":
						cmd := exec.Command("echo", "abcdefg")
						cmd.Env = os.Environ()
						return cmd
					case name == "git" && len(arg) > 0 && arg[0] == "status":
						cmd := exec.Command("echo", " M somefile.txt") // Dirty
						cmd.Env = os.Environ()
						return cmd
					}
					return exec.Command(name, arg...)
				}
			},
			expected: "dev-abcdefg-dirty",
		},
		{
			name: "Fallback to dev (all git commands fail)",
			setup: func() {
				execCommand = func(name string, arg ...string) *exec.Cmd {
					cmd := exec.Command("false") // Simulate all git commands failing
					cmd.Env = os.Environ()
					return cmd
				}
			},
			expected: "dev",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resetVersionState()                // Reset global state for each test
			originalExecCommand := execCommand // Store original to restore later
			defer func() {
				execCommand = originalExecCommand // Restore original after test
			}()

			tc.setup()
			got := resolveVersion()
			if got != tc.expected {
				t.Errorf("resolveVersion() got = %q, want %q", got, tc.expected)
			}
		})
	}
}

func TestGitDescribe(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	t.Run("gitDescribe_Success", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "v1.2.3")
		}
		tag, err := gitDescribe()
		if err != nil {
			t.Fatalf("gitDescribe failed: %v", err)
		}
		if tag != "v1.2.3" {
			t.Errorf("gitDescribe got %q, want %q", tag, "v1.2.3")
		}
	})

	t.Run("gitDescribe_Failure", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false") // Command that exits with non-zero
		}
		_, err := gitDescribe()
		if err == nil {
			t.Fatal("expected gitDescribe to fail")
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Errorf("expected *exec.ExitError, got %T", err)
		}
	})
}

func TestGitShortHash(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	t.Run("gitShortHash_Success", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "abc1234")
		}
		hash, err := gitShortHash()
		if err != nil {
			t.Fatalf("gitShortHash failed: %v", err)
		}
		if hash != "abc1234" {
			t.Errorf("gitShortHash got %q, want %q", hash, "abc1234")
		}
	})

	t.Run("gitShortHash_Failure", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false")
		}
		_, err := gitShortHash()
		if err == nil {
			t.Fatal("expected gitShortHash to fail")
		}
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) {
			t.Errorf("expected *exec.ExitError, got %T", err)
		}
	})
}

func TestGitDirty(t *testing.T) {
	originalExecCommand := execCommand
	defer func() { execCommand = originalExecCommand }()

	t.Run("gitDirty_Clean", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "") // Empty output means clean
		}
		if gitDirty() {
			t.Error("expected clean, got dirty")
		}
	})

	t.Run("gitDirty_Dirty", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", " M somefile.txt") // Non-empty output means dirty
		}
		if !gitDirty() {
			t.Error("expected dirty, got clean")
		}
	})

	t.Run("gitDirty_Failure", func(t *testing.T) {
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("false") // Simulate error
		}
		if gitDirty() { // Should return false (not dirty) on command error
			t.Error("expected clean on error, got dirty")
		}
	})
}
