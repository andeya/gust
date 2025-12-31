//go:build !windows
// +build !windows

package shutdown

import (
	"context"
	"os/signal"
	"sync/atomic"
	"syscall"
)

// Listen starts listening for shutdown signals (SIGINT, SIGTERM, SIGUSR2).
// This method blocks until a signal is received.
//
// On Unix systems:
//   - SIGINT, SIGTERM: Triggers graceful shutdown
//   - SIGUSR2: Triggers graceful reboot (if supported)
func (s *Shutdown) Listen() {
	// Fast path: check if already listening using atomic operation (no lock needed)
	// This avoids lock contention when Listen() is called frequently
	if atomic.LoadInt32(&s.listening) != 0 {
		return // already listening, fast path return
	}
	// Slow path: need to acquire lock
	s.mu.Lock()
	// Double-check after acquiring lock (another goroutine might have set it)
	if atomic.LoadInt32(&s.listening) != 0 {
		s.mu.Unlock()
		return // already listening
	}
	if atomic.LoadInt32(&s.stopped) != 0 {
		// If already stopped, create a new stopCh
		s.stopCh = make(chan struct{})
		atomic.StoreInt32(&s.stopped, 0)
	}
	atomic.StoreInt32(&s.listening, 1)
	s.mu.Unlock()

	// Stop any previous signal handlers to avoid accumulation
	signal.Stop(s.signalCh)
	// Clear any buffered signals (non-blocking drain)
	// Limit the number of iterations to prevent infinite loops
	maxDrainIterations := 1000 // Prevent infinite loop
	for i := 0; i < maxDrainIterations; i++ {
		select {
		case <-s.signalCh:
			// Drain buffered signals
		default:
			// No more buffered signals, break
			goto drainComplete
		}
	}
drainComplete:
	signal.Notify(s.signalCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)

	select {
	case sig := <-s.signalCh:
		signal.Stop(s.signalCh)
		s.mu.Lock()
		atomic.StoreInt32(&s.listening, 0)
		s.mu.Unlock()
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			s.Shutdown(context.Background())
		case syscall.SIGUSR2:
			s.Reboot(context.Background())
		}
	case <-s.stopCh:
		signal.Stop(s.signalCh)
		s.mu.Lock()
		atomic.StoreInt32(&s.listening, 0)
		s.mu.Unlock()
	}
}
