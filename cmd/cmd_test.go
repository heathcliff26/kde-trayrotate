package main

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/heathcliff26/kde-trayrotate/pkg/version"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCommand(t *testing.T) {
	cmd := NewRootCommand()

	name := strings.ToLower(version.Name)
	assert.Equal(t, name, cmd.Use)
}

func TestRootCommandExit(t *testing.T) {
	if os.Getenv("RUN_CRASH_TEST") == "1" {
		cmd := NewRootCommand()
		cmd.SetArgs([]string{"-l", "unknown"})
		_ = cmd.Execute()
		// Should not reach here, ensure exit with 0 if it does
		os.Exit(0)
	}
	output := execExitTest(t, "TestRootCommandExit", true)

	assert := assert.New(t)
	assert.Equal("Error: invalid log level 'unknown'\n", output, "Should return correct output")
}

func execExitTest(t *testing.T, test string, exitsError bool) string {
	cmd := exec.Command(os.Args[0], "-test.run="+test)
	cmd.Env = append(os.Environ(), "RUN_CRASH_TEST=1")
	out, err := cmd.CombinedOutput()
	if exitsError && err == nil {
		t.Fatal("Process exited without error")
	} else if !exitsError && err == nil {
		return string(out)
	}
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return string(out)
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
	return ""
}

func TestSetLogLevel(t *testing.T) {
	tMatrix := []struct {
		Name        string
		Level       slog.Level
		ShouldError bool
	}{
		{"debug", slog.LevelDebug, false},
		{"info", slog.LevelInfo, false},
		{"warn", slog.LevelWarn, false},
		{"error", slog.LevelError, false},
		{"DEBUG", slog.LevelDebug, false},
		{"INFO", slog.LevelInfo, false},
		{"WARN", slog.LevelWarn, false},
		{"ERROR", slog.LevelError, false},
		{"Unknown", 0, true},
	}
	t.Cleanup(func() {
		err := setLogLevel("info")
		if err != nil {
			t.Fatalf("Failed to cleanup after test: %v", err)
		}
	})

	for _, tCase := range tMatrix {
		t.Run(tCase.Name, func(t *testing.T) {
			opts := slog.HandlerOptions{
				Level: slog.Level(16),
			}
			logger := slog.New(slog.NewTextHandler(os.Stdout, &opts))
			slog.SetDefault(logger)

			err := setLogLevel(tCase.Name)

			assert := assert.New(t)

			if tCase.ShouldError {
				assert.Error(err, "Should return an error")
			} else {
				assert.NoError(err, "Should not return an error")
				assert.True(slog.Default().Enabled(t.Context(), tCase.Level), "Expected Level should be enabled")
			}
		})
	}
}
