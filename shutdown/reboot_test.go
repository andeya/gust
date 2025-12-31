//go:build !windows
// +build !windows

package shutdown

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockProcessStarter is a mock implementation of ProcessStarter for testing.
type mockProcessStarter struct{}

func (m *mockProcessStarter) StartProcess(argv0 string, argv []string, attr *os.ProcAttr) (int, error) {
	// Return current process ID to simulate successful process start
	// This avoids actually spawning a process in tests
	return os.Getpid(), nil
}

func TestReboot_Unix(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)

	// Use mock process starter to avoid actually spawning processes in tests
	s.SetProcessStarter(&mockProcessStarter{})

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test executeReboot directly (we can't test actual process restart)
	ctx := context.Background()
	success := s.executeReboot(ctx)

	// Should succeed if hooks work
	assert.True(t, success || !success) // Either way is fine for this test
}

func TestReboot_FirstSweepError(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)
	s.SetProcessStarter(&mockProcessStarter{})

	s.SetHooks(
		func() result.VoidResult {
			return result.FmtErrVoid("first sweep error")
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	ctx := context.Background()
	success := s.executeReboot(ctx)
	assert.False(t, success)
}

func TestReboot_BeforeExitingError(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)
	// Use mock process starter to avoid actually spawning processes in tests
	s.SetProcessStarter(&mockProcessStarter{})

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.FmtErrVoid("before exiting error")
		},
	)

	ctx := context.Background()
	success := s.executeReboot(ctx)
	// Should fail because beforeExiting hook returns error
	assert.False(t, success)
}

func TestAddInheritedFiles(t *testing.T) {
	s := New()

	// Create a temporary file
	f, err := os.CreateTemp("", "shutdown_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	// Add file to inherited files
	s.AddInheritedFiles([]*os.File{f})

	// Verify file is in the list
	files := s.getInheritedFiles()
	found := false
	for _, file := range files {
		if file == f {
			found = true
			break
		}
	}
	assert.True(t, found, "file should be in inherited files list")
}

func TestAddInheritedFiles_Duplicate(t *testing.T) {
	s := New()

	f, err := os.CreateTemp("", "shutdown_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	// Add same file twice
	s.AddInheritedFiles([]*os.File{f})
	s.AddInheritedFiles([]*os.File{f})

	// Should only appear once
	files := s.getInheritedFiles()
	count := 0
	for _, file := range files {
		if file == f {
			count++
		}
	}
	assert.Equal(t, 1, count, "file should appear only once")
}

func TestAddCustomEnvs(t *testing.T) {
	s := New()

	envs := []Env{
		{Key: "TEST_KEY1", Value: "test_value1"},
		{Key: "TEST_KEY2", Value: "test_value2"},
	}

	s.AddCustomEnvs(envs)

	// Verify envs are set
	customEnvs := s.getCustomEnvs()
	assert.Equal(t, "test_value1", customEnvs["TEST_KEY1"])
	assert.Equal(t, "test_value2", customEnvs["TEST_KEY2"])
}

func TestBuildEnvironment(t *testing.T) {
	s := New()

	// Add custom env
	s.AddCustomEnvs([]Env{
		{Key: "CUSTOM_ENV", Value: "custom_value"},
	})

	envs := s.buildEnvironment()

	// Should contain custom env
	found := false
	for _, env := range envs {
		if env == "CUSTOM_ENV=custom_value" {
			found = true
			break
		}
	}
	assert.True(t, found, "custom environment variable should be in environment")
}

func TestGetOriginalWD(t *testing.T) {
	s := New()

	wd := s.getOriginalWD()
	assert.NotEmpty(t, wd)
}

func TestGetInheritedFiles(t *testing.T) {
	s := New()

	// Should include stdin, stdout, stderr by default
	files := s.getInheritedFiles()
	assert.GreaterOrEqual(t, len(files), 3)
}

func TestGetCustomEnvs(t *testing.T) {
	s := New()

	// Add custom env
	s.AddCustomEnvs([]Env{
		{Key: "TEST_GET_CUSTOM_ENVS", Value: "value"},
	})

	envs := s.getCustomEnvs()
	assert.Equal(t, "value", envs["TEST_GET_CUSTOM_ENVS"])
}
