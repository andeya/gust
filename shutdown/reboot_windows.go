//go:build windows
// +build windows

package shutdown

import (
	"context"
	"os"
)

// Reboot gracefully reboots the application.
//
// NOTE: Reboot is not supported on Windows.
// This method logs a warning and calls Shutdown instead.
func (s *Shutdown) Reboot(ctx context.Context) {
	if ctx == nil {
		ctx = context.Background()
	}
	s.logf(ctx, "Windows system doesn't support reboot! Calling Shutdown() instead.")
	s.Shutdown(ctx)
}

// Env represents an environment variable to be inherited by the new process.
// On Windows, this is a no-op.
type Env struct {
	Key   string
	Value string
}

// AddInheritedFiles adds files to be inherited by the new process during reboot.
// On Windows, this is a no-op.
func (s *Shutdown) AddInheritedFiles(files []*os.File) {
	// No-op on Windows
}

// AddCustomEnvs adds custom environment variables to be inherited by the new process.
// On Windows, this is a no-op.
func (s *Shutdown) AddCustomEnvs(envs []Env) {
	// No-op on Windows
}
