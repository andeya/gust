package shutdown

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// TestCPUUsage_StopMultipleTimes tests that calling Stop() multiple times doesn't cause issues
func TestCPUUsage_StopMultipleTimes(t *testing.T) {
	s := New()

	// Start listening
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	// Call Stop() multiple times - should not panic or cause high CPU
	for i := 0; i < 10; i++ {
		s.Stop()
		time.Sleep(10 * time.Millisecond)
	}

	// Wait for Listen to return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after Stop()")
	}
}

// TestCPUUsage_ListenAfterStop tests that Listen() can be called again after Stop()
func TestCPUUsage_ListenAfterStop(t *testing.T) {
	s := New()

	// First Listen
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

	// Second Listen after Stop
	done2 := make(chan bool, 1)
	go func() {
		s.Listen()
		done2 <- true
	}()

	time.Sleep(50 * time.Millisecond)
	s.Stop()

	select {
	case <-done2:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Second Listen did not return")
	}
}

// TestCPUUsage_ConcurrentStop tests concurrent Stop() calls
func TestCPUUsage_ConcurrentStop(t *testing.T) {
	s := New()

	// Start listening
	done := make(chan bool, 1)
	go func() {
		s.Listen()
		done <- true
	}()

	time.Sleep(50 * time.Millisecond)

	// Concurrent Stop() calls
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Stop()
		}()
	}
	wg.Wait()

	// Should not panic and Listen should return
	select {
	case <-done:
		// Success
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after concurrent Stop()")
	}
}

// TestCPUUsage_NoBusyWait tests that Listen() doesn't busy wait
func TestCPUUsage_NoBusyWait(t *testing.T) {
	s := New()

	// Start listening - should block without consuming CPU
	done := make(chan bool, 1)
	start := time.Now()
	go func() {
		s.Listen()
		done <- true
	}()

	// Wait a bit - Listen should be blocking, not consuming CPU
	time.Sleep(100 * time.Millisecond)

	// Stop it
	s.Stop()

	select {
	case <-done:
		elapsed := time.Since(start)
		// Should have taken at least 100ms (the sleep time)
		assert.True(t, elapsed >= 100*time.Millisecond, "Listen should have blocked")
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

// TestCPUUsage_MultipleListeners tests multiple Listen() calls
func TestCPUUsage_MultipleListeners(t *testing.T) {
	s := New()

	// Multiple Listen() calls - only first should actually listen
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.Listen() // Should return immediately if already listening
		}()
	}

	// Give them time to start
	time.Sleep(50 * time.Millisecond)

	// Only one should be listening
	assert.True(t, s.IsListening())

	// Stop
	s.Stop()

	// Wait for all Listen() calls to return (with timeout)
	done := make(chan bool, 1)
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		// Success
	case <-time.After(2 * time.Second):
		t.Fatal("Not all Listen() calls returned within timeout")
	}
}
