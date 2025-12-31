//go:build !windows
// +build !windows

package shutdown

import (
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
