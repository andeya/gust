//go:build !windows
// +build !windows

package shutdown

import (
	"context"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	// Mock dependencies to prevent actual process operations during tests
	globalDeps.startProcess = func(argv0 string, argv []string, attr *os.ProcAttr) (int, error) {
		return os.Getpid(), nil
	}
	globalDeps.killProcess = func(pid int, sig syscall.Signal) error {
		return nil
	}
	globalDeps.exit = func(code int) {}
}

func TestExecuteReboot(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	// Success case
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)
	assert.True(t, s.executeReboot(context.Background()))

	// FirstSweep error
	s.SetHooks(
		func() result.VoidResult { return result.FmtErrVoid("error") },
		func() result.VoidResult { return result.OkVoid() },
	)
	assert.False(t, s.executeReboot(context.Background()))

	// BeforeExiting error
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.FmtErrVoid("error") },
	)
	assert.False(t, s.executeReboot(context.Background()))

	// Panic recovery
	s.SetHooks(
		func() result.VoidResult { panic("test panic") },
		func() result.VoidResult { return result.OkVoid() },
	)
	assert.False(t, s.executeReboot(context.Background()))
}

func TestRebootInternal(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	// With nil context
	//nolint:staticcheck // intentionally pass nil
	assert.True(t, s.rebootInternal(nil))

	// With context
	assert.True(t, s.rebootInternal(context.Background()))

	// Timeout path
	s.SetHooks(
		func() result.VoidResult { time.Sleep(100 * time.Millisecond); return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	assert.False(t, s.rebootInternal(ctx))
}

func TestReboot(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)
	var called bool
	s.SetHooks(
		func() result.VoidResult { called = true; return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)
	s.Reboot(context.Background())
	assert.True(t, called)
}

func TestInheritedFiles(t *testing.T) {
	s := New()

	// Default files (stdin, stdout, stderr)
	files := s.getInheritedFiles()
	assert.GreaterOrEqual(t, len(files), 3)

	// Add file
	f, err := os.CreateTemp("", "test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	s.AddInheritedFiles([]*os.File{f})
	files = s.getInheritedFiles()
	assert.Contains(t, files, f)

	// Add duplicate (should not duplicate)
	beforeLen := len(s.getInheritedFiles())
	s.AddInheritedFiles([]*os.File{f})
	assert.Equal(t, beforeLen, len(s.getInheritedFiles()))
}

func TestCustomEnvs(t *testing.T) {
	s := New()

	s.AddCustomEnvs([]Env{{Key: "K1", Value: "V1"}, {Key: "K2", Value: "V2"}})
	envs := s.getCustomEnvs()
	assert.Equal(t, "V1", envs["K1"])
	assert.Equal(t, "V2", envs["K2"])
}

func TestBuildEnvironment(t *testing.T) {
	s := New()

	// Custom env
	s.AddCustomEnvs([]Env{{Key: "CUSTOM_KEY", Value: "custom_value"}})
	envs := s.buildEnvironment()
	assert.Contains(t, envs, "CUSTOM_KEY=custom_value")

	// Overwrite system env
	os.Setenv("TEST_OVERWRITE", "system")
	defer os.Unsetenv("TEST_OVERWRITE")
	s.AddCustomEnvs([]Env{{Key: "TEST_OVERWRITE", Value: "custom"}})
	envs = s.buildEnvironment()
	assert.Contains(t, envs, "TEST_OVERWRITE=custom")
	for _, e := range envs {
		assert.NotEqual(t, "TEST_OVERWRITE=system", e)
	}
}

func TestStartProcess(t *testing.T) {
	s := New()

	// Success
	res := s.startProcess()
	assert.True(t, res.IsOk())

	// LookPath error
	orig := os.Args
	os.Args = []string{"/nonexistent/path"}
	defer func() { os.Args = orig }()
	res = s.startProcess()
	assert.True(t, res.IsErr())
}

func TestKillProcess(t *testing.T) {
	s := New()
	assert.NoError(t, s.killProcess(os.Getpid(), syscall.SIGTERM))

	// Failure case
	orig := globalDeps.killProcess
	globalDeps.killProcess = func(pid int, sig syscall.Signal) error { return syscall.EPERM }
	defer func() { globalDeps.killProcess = orig }()

	if os.Getppid() != 1 {
		assert.False(t, s.rebootInternal(context.Background()))
	}
}

func TestGetOriginalWD(t *testing.T) {
	s := New()
	assert.NotEmpty(t, s.getOriginalWD())
}
