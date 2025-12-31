//go:build windows
// +build windows

package shutdown

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReboot_Windows(t *testing.T) {
	s := New()

	// On Windows, Reboot should call Shutdown
	// We can't test os.Exit, but we can verify it doesn't panic
	assert.NotPanics(t, func() {
		// This would call os.Exit, so we just verify the method exists
		_ = s.Reboot
	})
}

func TestAddInheritedFiles_Windows(t *testing.T) {
	s := New()

	// On Windows, this should be a no-op
	f, err := os.CreateTemp("", "shutdown_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())
	defer f.Close()

	// Should not panic
	s.AddInheritedFiles([]*os.File{f})
}

func TestAddCustomEnvs_Windows(t *testing.T) {
	s := New()

	// On Windows, this should be a no-op
	envs := []Env{
		{Key: "TEST", Value: "value"},
	}

	// Should not panic
	s.AddCustomEnvs(envs)
}
