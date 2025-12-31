//go:build !windows
// +build !windows

package shutdown

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"syscall"

	"github.com/andeya/gust/result"
)

// Env represents an environment variable to be inherited by the new process.
type Env struct {
	Key   string
	Value string
}

// rebootInternal contains the internal reboot logic without os.Exit.
// This allows testing the logic without terminating the test process.
// Returns true if reboot was graceful, false otherwise.
func (s *Shutdown) rebootInternal(ctx context.Context) bool {
	// Use provided context or create one with timeout
	rebootCtx := ctx
	if rebootCtx == nil {
		var cancel context.CancelFunc
		rebootCtx, cancel = context.WithTimeout(context.Background(), s.Timeout())
		defer cancel()
	}

	s.logf(rebootCtx, "rebooting process...")

	ppid := os.Getppid()
	graceful := true

	// Execute reboot in a goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		graceful = s.executeReboot(rebootCtx)
	}()

	// Wait for completion or timeout
	select {
	case <-rebootCtx.Done():
		if err := rebootCtx.Err(); err != nil {
			s.logErrorf(rebootCtx, "reboot timeout: %v", err)
		}
		graceful = false
	case <-done:
		// Reboot completed
	}

	// Kill parent process if needed
	if ppid != 1 {
		if err := s.killProcess(ppid, syscall.SIGTERM); err != nil {
			s.logErrorf(rebootCtx, "failed to kill parent process: %v", err)
			graceful = false
		}
	}

	if graceful {
		s.logf(rebootCtx, "process rebooted gracefully")
	} else {
		s.logf(rebootCtx, "process rebooted, but not gracefully")
	}
	return graceful
}

// Reboot gracefully reboots the application by starting a new process
// and then shutting down the current one.
//
// NOTE: Reboot is not supported on Windows. On Windows, this method
// will log a warning and exit.
func (s *Shutdown) Reboot(ctx context.Context) {
	graceful := s.rebootInternal(ctx)
	if graceful {
		globalDeps.exit(0)
	} else {
		globalDeps.exit(-1)
	}
}

// executeReboot executes the reboot sequence.
func (s *Shutdown) executeReboot(ctx context.Context) (r bool) {
	defer func() {
		if p := recover(); p != nil {
			s.logErrorf(ctx, "panic during reboot: %v", p)
			r = false
		}
	}()

	// Execute first sweep
	s.logf(ctx, "executing first sweep...")
	firstSweepRes := s.getFirstSweep()()
	if firstSweepRes.IsErr() {
		s.logErrorf(ctx, "first sweep failed: %v", firstSweepRes.Err())
		return false
	}

	// Start new process
	s.logf(ctx, "starting new process...")
	startRes := s.startProcess()
	if startRes.IsErr() {
		s.logErrorf(ctx, "failed to start new process: %v", startRes.Err())
		return false
	}

	// Execute before exiting
	s.logf(ctx, "executing before exiting...")
	beforeExitingRes := s.getBeforeExiting()()
	if beforeExitingRes.IsErr() {
		s.logErrorf(ctx, "before exiting failed: %v", beforeExitingRes.Err())
		return false
	}

	return true
}

// startProcess starts a new process with the same arguments and environment.
func (s *Shutdown) startProcess() (r result.Result[int]) {
	defer r.Catch()

	// Get executable path
	argv0 := result.Ret(exec.LookPath(os.Args[0])).Unwrap()

	// Build environment
	envs := s.buildEnvironment()

	// Get inherited files
	files := s.getInheritedFiles()

	pid, err := globalDeps.startProcess(argv0, os.Args, &os.ProcAttr{
		Dir:   s.getOriginalWD(),
		Env:   envs,
		Files: files,
	})
	if err != nil {
		return result.TryErr[int](err)
	}
	return result.Ok(pid)
}

// buildEnvironment builds the environment variables for the new process.
func (s *Shutdown) buildEnvironment() []string {
	// Get custom environments
	customEnvs := s.getCustomEnvs()

	// Start with system environment
	envs := make([]string, 0, len(os.Environ())+len(customEnvs))
	for _, env := range os.Environ() {
		k := strings.Split(env, "=")[0]
		if _, ok := customEnvs[k]; !ok {
			envs = append(envs, env)
		}
	}

	// Add custom environments
	for k, v := range customEnvs {
		envs = append(envs, k+"="+v)
	}

	return envs
}

// InheritedFiles manages files to be inherited by the new process.
type InheritedFiles struct {
	mu    sync.Mutex
	files []*os.File
}

var defaultInheritedFiles = &InheritedFiles{
	files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
}

// AddInheritedFiles adds files to be inherited by the new process during reboot.
func (s *Shutdown) AddInheritedFiles(files []*os.File) {
	defaultInheritedFiles.mu.Lock()
	defer defaultInheritedFiles.mu.Unlock()

	for _, f := range files {
		// Check if already added
		exists := false
		for _, existing := range defaultInheritedFiles.files {
			if existing == f {
				exists = true
				break
			}
		}
		if !exists {
			defaultInheritedFiles.files = append(defaultInheritedFiles.files, f)
		}
	}
}

// getInheritedFiles returns the files to be inherited.
func (s *Shutdown) getInheritedFiles() []*os.File {
	defaultInheritedFiles.mu.Lock()
	defer defaultInheritedFiles.mu.Unlock()

	// Return a copy to avoid modification
	files := make([]*os.File, len(defaultInheritedFiles.files))
	copy(files, defaultInheritedFiles.files)
	return files
}

// CustomEnvs manages custom environment variables.
type CustomEnvs struct {
	mu   sync.Mutex
	envs map[string]string
}

var defaultCustomEnvs = &CustomEnvs{
	envs: make(map[string]string),
}

// AddCustomEnvs adds custom environment variables to be inherited by the new process.
func (s *Shutdown) AddCustomEnvs(envs []Env) {
	defaultCustomEnvs.mu.Lock()
	defer defaultCustomEnvs.mu.Unlock()

	for _, env := range envs {
		defaultCustomEnvs.envs[env.Key] = env.Value
	}
}

// getCustomEnvs returns the custom environment variables.
func (s *Shutdown) getCustomEnvs() map[string]string {
	defaultCustomEnvs.mu.Lock()
	defer defaultCustomEnvs.mu.Unlock()

	// Return a copy to avoid modification
	envs := make(map[string]string, len(defaultCustomEnvs.envs))
	for k, v := range defaultCustomEnvs.envs {
		envs[k] = v
	}
	return envs
}

// originalWD stores the original working directory.
var originalWD = func() string {
	wd, _ := os.Getwd()
	return wd
}()

// getOriginalWD returns the original working directory.
func (s *Shutdown) getOriginalWD() string {
	return originalWD
}

// killProcess kills a process using the configured killer.
func (s *Shutdown) killProcess(pid int, sig syscall.Signal) error {
	return globalDeps.killProcess(pid, sig)
}
