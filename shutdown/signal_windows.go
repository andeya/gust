//go:build windows
// +build windows

package shutdown

import (
	"context"
	"os"
	"os/signal"
	"sync/atomic"
)

// Listen starts listening for shutdown signals (SIGINT, SIGKILL).
// This method blocks until a signal is received.
//
// On Windows:
//   - SIGINT, SIGKILL: Triggers graceful shutdown
//   - Reboot is not supported on Windows
func (s *Shutdown) Listen() {
	// Fast path: check if already listening using atomic operation (no lock needed)
	if atomic.LoadInt32(&s.listening) != 0 {
		return // already listening, fast path return
	}
	s.mu.Lock()
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
	// On Windows, only SIGINT can be trapped; SIGKILL cannot be trapped
	signal.Notify(s.signalCh, os.Interrupt)

	select {
	case <-s.signalCh:
		signal.Stop(s.signalCh)
		s.mu.Lock()
		atomic.StoreInt32(&s.listening, 0)
		s.mu.Unlock()
		s.Shutdown(context.Background())
	case <-s.stopCh:
		signal.Stop(s.signalCh)
		s.mu.Lock()
		atomic.StoreInt32(&s.listening, 0)
		s.mu.Unlock()
	}
}
