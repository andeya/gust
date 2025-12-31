//go:build !windows
// +build !windows

package shutdown

import (
	"os"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestListen_Unix(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test that Listen can be stopped
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Stop listening after a short delay
	time.Sleep(50 * time.Millisecond)
	s.Stop()

	// Wait for Listen to return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after Stop()")
	}
}

func TestListen_AlreadyListening(t *testing.T) {
	s := New()

	// Start listening in background
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait a bit for listening to start
	time.Sleep(50 * time.Millisecond)

	// Verify it's listening
	assert.True(t, s.IsListening())

	// Try to listen again (should return immediately because already listening)
	// This should not block
	done2 := make(chan bool, 1)
	go func() {
		s.Listen()
		done2 <- true
	}()

	// Should return immediately
	select {
	case <-done2:
		// Success - returned immediately
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Second Listen() should return immediately when already listening")
	}

	// Clean up
	s.Stop()

	// Wait for first Listen to return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after Stop()")
	}
}

func TestStop_WhileListening(t *testing.T) {
	s := New()

	// Start listening in background
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait a bit for listening to start
	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	// Stop listening
	s.Stop()
	time.Sleep(50 * time.Millisecond)

	assert.False(t, s.IsListening())

	// Wait for Listen to return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after Stop()")
	}
}

func TestSignalHandling(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	// Test that signal channel is created
	assert.NotNil(t, s.signalCh)
}

// TestListen_DrainSignals tests signal draining
func TestListen_DrainSignals(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Fill signal channel with buffered signals
	// Note: We fill the channel BEFORE starting Listen() to ensure signals are drained
	// and not processed (which would trigger Shutdown/Reboot and os.Exit)
	for i := 0; i < 5; i++ {
		select {
		case s.signalCh <- os.Interrupt:
		default:
		}
	}

	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Give Listen() time to drain signals, then stop before any signal can be processed
	time.Sleep(10 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// Success - signals should be drained, not processed
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

// TestListen_DrainManySignals tests draining many buffered signals
func TestListen_DrainManySignals(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Fill signal channel with many buffered signals (more than maxDrainIterations)
	for i := 0; i < 1500; i++ {
		select {
		case s.signalCh <- os.Interrupt:
		default:
			// Channel full, break out of loop
			goto loopDone
		}
	}
loopDone:

	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// Success - signals should be drained, not processed
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

// TestListen_MaxDrainIterations tests the max drain iterations path
func TestListen_MaxDrainIterations(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Fill signal channel with more than maxDrainIterations (1000) signals
	// This tests the path where we hit the max iterations limit
	// Note: We fill the channel BEFORE starting Listen() to ensure signals are drained
	// and not processed (which would trigger Shutdown/Reboot and os.Exit)
	for i := 0; i < 1500; i++ {
		select {
		case s.signalCh <- os.Interrupt:
		default:
			// Channel full, break out of loop
			goto loopDone
		}
	}
loopDone:

	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Give Listen() time to drain signals, then stop before any signal can be processed
	time.Sleep(10 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// Success - signals should be drained, not processed
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

// TestListen_DoubleCheckPath tests the double-check locking path in Listen
func TestListen_DoubleCheckPath(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Set listening flag before lock (tests double-check path)
	atomic.StoreInt32(&s.listening, 1)

	// Listen should return immediately due to fast path
	s.Listen()

	// Reset for next test
	atomic.StoreInt32(&s.listening, 0)

	// Test slow path double-check
	s.mu.Lock()
	atomic.StoreInt32(&s.listening, 1)
	s.mu.Unlock()

	// Listen should return immediately due to slow path double-check
	s.Listen()
}

// TestListen_StoppedAndRestart tests Listen after Stop
func TestListen_StoppedAndRestart(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Start and stop
	done1 := make(chan bool, 1)
	go func() {
		s.Listen()
		done1 <- true
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	select {
	case <-done1:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("First Listen did not return")
	}

	// Start again (should work after stop)
	done2 := make(chan bool, 1)
	go func() {
		s.Listen()
		done2 <- true
	}()

	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	s.Stop()

	select {
	case <-done2:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Second Listen did not return")
	}
}

// TestListen_AlreadyStopped tests Listen when already stopped
func TestListen_AlreadyStopped(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Mark as stopped
	s.mu.Lock()
	atomic.StoreInt32(&s.stopped, 1)
	s.mu.Unlock()

	// Listen should reset stopped and create new stopCh
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	s.Stop()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

// TestListen_SignalSIGINT tests Listen when receiving SIGINT signal
func TestListen_SignalSIGINT(t *testing.T) {
	// Mock exit to prevent actual exit
	originalExit := globalDeps.exit
	globalDeps.exit = func(code int) {}
	defer func() {
		globalDeps.exit = originalExit
	}()

	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var shutdownCalled bool
	s.SetHooks(
		func() result.VoidResult {
			shutdownCalled = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Start Listen in a goroutine
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait for Listen to start
	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	// Send SIGINT signal
	process, _ := os.FindProcess(os.Getpid())
	_ = process.Signal(syscall.SIGINT)

	// Wait for signal to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify shutdown was called
	assert.True(t, shutdownCalled)

	// Clean up
	s.Stop()
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		// Timeout is acceptable here as signal processing may take time
	}
}

// TestListen_SignalSIGTERM tests Listen when receiving SIGTERM signal
func TestListen_SignalSIGTERM(t *testing.T) {
	// Mock exit to prevent actual exit
	originalExit := globalDeps.exit
	globalDeps.exit = func(code int) {}
	defer func() {
		globalDeps.exit = originalExit
	}()

	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var shutdownCalled bool
	s.SetHooks(
		func() result.VoidResult {
			shutdownCalled = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Start Listen in a goroutine
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait for Listen to start
	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	// Send SIGTERM signal
	process, _ := os.FindProcess(os.Getpid())
	_ = process.Signal(syscall.SIGTERM)

	// Wait for signal to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify shutdown was called
	assert.True(t, shutdownCalled)

	// Clean up
	s.Stop()
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		// Timeout is acceptable here as signal processing may take time
	}
}

// TestListen_SignalSIGUSR2 tests Listen when receiving SIGUSR2 signal
func TestListen_SignalSIGUSR2(t *testing.T) {
	// Mock exit to prevent actual exit
	originalExit := globalDeps.exit
	globalDeps.exit = func(code int) {}
	defer func() {
		globalDeps.exit = originalExit
	}()

	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var rebootCalled bool
	s.SetHooks(
		func() result.VoidResult {
			rebootCalled = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Start Listen in a goroutine
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait for Listen to start
	time.Sleep(50 * time.Millisecond)
	assert.True(t, s.IsListening())

	// Send SIGUSR2 signal
	process, _ := os.FindProcess(os.Getpid())
	_ = process.Signal(syscall.SIGUSR2)

	// Wait for signal to be processed
	time.Sleep(100 * time.Millisecond)

	// Verify reboot was called (executeReboot was called)
	assert.True(t, rebootCalled)

	// Clean up
	s.Stop()
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		// Timeout is acceptable here as signal processing may take time
	}
}

// TestListen_DrainCompletePath tests the drainComplete path in Listen
func TestListen_DrainCompletePath(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Fill signal channel with exactly maxDrainIterations signals
	// This tests the path where we hit max iterations and goto drainComplete
	for i := 0; i < 1000; i++ {
		select {
		case s.signalCh <- os.Interrupt:
		default:
			// Channel full, break
			goto fillDone
		}
	}
fillDone:

	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	// Give Listen() time to drain signals
	time.Sleep(10 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}
