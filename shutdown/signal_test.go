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

func init() {
	globalDeps.exit = func(code int) {}
}

func TestListen(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	// Test Listen and Stop
	done := make(chan bool, 1)
	go func() { s.Listen(); done <- true }()
	time.Sleep(30 * time.Millisecond)
	assert.True(t, s.IsListening())
	s.Stop()
	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return after Stop")
	}
}

func TestListenAlreadyListening(t *testing.T) {
	s := New()
	done := make(chan bool, 1)
	go func() { s.Listen(); done <- true }()
	time.Sleep(30 * time.Millisecond)

	// Second Listen should return immediately
	done2 := make(chan bool, 1)
	go func() { s.Listen(); done2 <- true }()
	select {
	case <-done2:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Second Listen should return immediately")
	}

	s.Stop()
	<-done
}

func TestListenStoppedAndRestart(t *testing.T) {
	s := New()
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	// First listen and stop
	done1 := make(chan bool, 1)
	go func() { s.Listen(); done1 <- true }()
	time.Sleep(30 * time.Millisecond)
	s.Stop()
	<-done1

	// Listen again after stop
	done2 := make(chan bool, 1)
	go func() { s.Listen(); done2 <- true }()
	time.Sleep(30 * time.Millisecond)
	assert.True(t, s.IsListening())
	s.Stop()
	<-done2
}

func TestListenSignals(t *testing.T) {
	tests := []struct {
		name   string
		signal syscall.Signal
	}{
		{"SIGINT", syscall.SIGINT},
		{"SIGTERM", syscall.SIGTERM},
		{"SIGUSR2", syscall.SIGUSR2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New()
			s.SetTimeout(50 * time.Millisecond)
			var called bool
			s.SetHooks(
				func() result.VoidResult { called = true; return result.OkVoid() },
				func() result.VoidResult { return result.OkVoid() },
			)

			done := make(chan bool, 1)
			go func() { s.Listen(); done <- true }()
			time.Sleep(30 * time.Millisecond)

			process, _ := os.FindProcess(os.Getpid())
			_ = process.Signal(tt.signal)
			time.Sleep(100 * time.Millisecond)

			assert.True(t, called)
			s.Stop()
			select {
			case <-done:
			case <-time.After(1 * time.Second):
			}
		})
	}
}

func TestListenDrainSignals(t *testing.T) {
	s := New()
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	// Fill signal channel before Listen
	for i := 0; i < 5; i++ {
		select {
		case s.signalCh <- os.Interrupt:
		default:
		}
	}

	done := make(chan bool, 1)
	go func() { s.Listen(); done <- true }()
	time.Sleep(10 * time.Millisecond)
	s.Stop()

	select {
	case <-done:
	case <-time.After(1 * time.Second):
		t.Fatal("Listen did not return")
	}
}

func TestListenDoubleCheck(t *testing.T) {
	s := New()

	// Fast path: already listening
	atomic.StoreInt32(&s.listening, 1)
	s.Listen()
	atomic.StoreInt32(&s.listening, 0)

	// Slow path: set listening during lock
	s.mu.Lock()
	atomic.StoreInt32(&s.listening, 1)
	s.mu.Unlock()
	s.Listen()
}

func TestListenAlreadyStopped(t *testing.T) {
	s := New()
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	atomic.StoreInt32(&s.stopped, 1)

	done := make(chan bool, 1)
	go func() { s.Listen(); done <- true }()
	time.Sleep(30 * time.Millisecond)
	assert.True(t, s.IsListening())
	s.Stop()
	<-done
}
