// Package shutdown provides graceful shutdown and reboot functionality for Go applications.
//
// This package allows applications to gracefully shut down or reboot by handling
// system signals (SIGINT, SIGTERM, SIGUSR2) and executing cleanup hooks before exiting.
//
// # Examples
//
//	```go
//	// Create a shutdown manager
//	shutdown := shutdown.New()
//
//	// Set hooks
//	shutdown.SetHooks(
//	    func() result.VoidResult {
//	        // First sweep: close connections, stop accepting new requests
//	        return result.OkVoid()
//	    },
//	    func() result.VoidResult {
//	        // Before exiting: final cleanup
//	        return result.OkVoid()
//	    },
//	)
//
//	// Start listening for signals
//	shutdown.Listen()
//
//	// Or manually trigger shutdown
//	shutdown.Shutdown(context.Background())
//	```
package shutdown

import (
	"context"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/andeya/gust/result"
)

// MinShutdownTimeout is the default minimum timeout for graceful shutdown.
const MinShutdownTimeout = 15 * time.Second

// processStarterFunc is a function type for starting new processes.
type processStarterFunc func(argv0 string, argv []string, attr *os.ProcAttr) (int, error)

// processKillerFunc is a function type for killing processes.
type processKillerFunc func(pid int, sig syscall.Signal) error

// osExitFunc is a function type for exiting the process.
type osExitFunc func(code int)

// globalDeps is the global dependency object used by all Shutdown instances.
// In tests, this can be replaced with a mock implementation.
var globalDeps = struct {
	startProcess processStarterFunc
	killProcess  processKillerFunc
	exit         osExitFunc
}{
	startProcess: func(argv0 string, argv []string, attr *os.ProcAttr) (int, error) {
		process, err := os.StartProcess(argv0, argv, attr)
		if err != nil {
			return 0, err
		}
		return process.Pid, nil
	},
	killProcess: syscall.Kill,
	exit:        os.Exit,
}

// Shutdown manages graceful shutdown and reboot of the application.
type Shutdown struct {
	timeout       time.Duration
	firstSweep    func() result.VoidResult
	beforeExiting func() result.VoidResult
	logger        Logger
	mu            sync.Mutex
	signalCh      chan os.Signal
	stopCh        chan struct{}
	listening     int32 // Use atomic for fast path check
	stopped       int32 // Use atomic for fast path check
}

// Logger is an interface for logging shutdown events.
// If nil, no logging is performed.
type Logger interface {
	Infof(ctx context.Context, format string, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
}

// New creates a new Shutdown instance with default settings.
func New() *Shutdown {
	return &Shutdown{
		timeout:       MinShutdownTimeout,
		firstSweep:    func() result.VoidResult { return result.OkVoid() },
		beforeExiting: func() result.VoidResult { return result.OkVoid() },
		signalCh:      make(chan os.Signal, 1),
		stopCh:        make(chan struct{}),
		listening:     0,
		stopped:       0,
	}
}

// SetTimeout sets the timeout for graceful shutdown.
// If timeout < 0, it's set to a very large value (effectively infinite).
// If 0 <= timeout < MinShutdownTimeout, it's set to MinShutdownTimeout.
func (s *Shutdown) SetTimeout(timeout time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if timeout < 0 {
		s.timeout = 1<<63 - 1 // effectively infinite
	} else if timeout < MinShutdownTimeout {
		s.timeout = MinShutdownTimeout
	} else {
		s.timeout = timeout
	}
}

// Timeout returns the current shutdown timeout.
func (s *Shutdown) Timeout() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.timeout
}

// SetHooks sets the hooks to be executed during shutdown.
//   - firstSweep: Executed first (e.g., close connections, stop accepting new requests)
//   - beforeExiting: Executed before process exits (e.g., final cleanup)
//
// If a hook is nil, it's replaced with a no-op function.
func (s *Shutdown) SetHooks(firstSweep, beforeExiting func() result.VoidResult) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if firstSweep == nil {
		firstSweep = func() result.VoidResult { return result.OkVoid() }
	}
	if beforeExiting == nil {
		beforeExiting = func() result.VoidResult { return result.OkVoid() }
	}

	s.firstSweep = firstSweep
	s.beforeExiting = beforeExiting
}

// SetLogger sets the logger for shutdown events.
// If nil, no logging is performed.
func (s *Shutdown) SetLogger(logger Logger) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logger = logger
}

// Logger returns the current logger.
func (s *Shutdown) Logger() Logger {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.logger
}

// shutdownInternal contains the internal shutdown logic without os.Exit.
// This allows testing the logic without terminating the test process.
func (s *Shutdown) shutdownInternal(ctx context.Context) {
	// Use provided context or create one with timeout
	shutdownCtx := ctx
	if shutdownCtx == nil {
		var cancel context.CancelFunc
		shutdownCtx, cancel = context.WithTimeout(context.Background(), s.Timeout())
		defer cancel()
	}

	s.logf(shutdownCtx, "shutting down process...")

	// Execute shutdown in a goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.executeShutdown(shutdownCtx)
	}()

	// Wait for completion or timeout
	select {
	case <-shutdownCtx.Done():
		if err := shutdownCtx.Err(); err != nil {
			s.logf(shutdownCtx, "shutdown timeout: %v", err)
		}
	case <-done:
		// Shutdown completed
	}

	s.logf(shutdownCtx, "process shutdown complete")
}

// Shutdown gracefully shuts down the application.
// It executes the hooks within the specified timeout and then exits with code 0.
//
// If ctx is provided, it's used instead of the configured timeout.
// The function will block until shutdown is complete or timeout occurs.
func (s *Shutdown) Shutdown(ctx context.Context) {
	defer globalDeps.exit(0)
	s.shutdownInternal(ctx)
}

// executeShutdown executes the shutdown hooks.
func (s *Shutdown) executeShutdown(ctx context.Context) (r result.VoidResult) {
	defer r.Catch()

	// Execute first sweep
	s.logf(ctx, "executing first sweep...")
	firstSweepRes := s.getFirstSweep()()
	if firstSweepRes.IsErr() {
		s.logf(ctx, "first sweep failed: %v", firstSweepRes.Err())
		return firstSweepRes
	}

	// Execute before exiting
	s.logf(ctx, "executing before exiting...")
	beforeExitingRes := s.getBeforeExiting()()
	if beforeExitingRes.IsErr() {
		s.logf(ctx, "before exiting failed: %v", beforeExitingRes.Err())
		return beforeExitingRes
	}

	s.logf(ctx, "process shut down gracefully")
	return result.OkVoid()
}

// getFirstSweep returns the first sweep hook (thread-safe).
func (s *Shutdown) getFirstSweep() func() result.VoidResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.firstSweep
}

// getBeforeExiting returns the before exiting hook (thread-safe).
func (s *Shutdown) getBeforeExiting() func() result.VoidResult {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.beforeExiting
}

// logf logs a message if logger is set.
func (s *Shutdown) logf(ctx context.Context, format string, args ...interface{}) {
	logger := s.Logger()
	if logger != nil {
		logger.Infof(ctx, format, args...)
	}
}

// logErrorf logs an error message if logger is set.
func (s *Shutdown) logErrorf(ctx context.Context, format string, args ...interface{}) {
	logger := s.Logger()
	if logger != nil {
		logger.Errorf(ctx, format, args...)
	}
}

// IsListening returns whether the shutdown manager is currently listening for signals.
func (s *Shutdown) IsListening() bool {
	return atomic.LoadInt32(&s.listening) != 0
}

// Stop stops listening for signals.
func (s *Shutdown) Stop() {
	// Fast path: check if already stopped using atomic operation (no lock needed)
	// This avoids lock contention when Stop() is called frequently
	if atomic.LoadInt32(&s.stopped) != 0 {
		return // already stopped, fast path return
	}
	// Fast path: check if not listening (no need to stop)
	if atomic.LoadInt32(&s.listening) == 0 {
		return // not listening, nothing to stop
	}
	// Slow path: need to acquire lock
	s.mu.Lock()
	defer s.mu.Unlock()

	// Double-check after acquiring lock (another goroutine might have changed state)
	if atomic.LoadInt32(&s.listening) != 0 && atomic.LoadInt32(&s.stopped) == 0 {
		close(s.stopCh)
		atomic.StoreInt32(&s.listening, 0)
		atomic.StoreInt32(&s.stopped, 1)
	}
}
