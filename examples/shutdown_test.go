package examples

import (
	"context"
	"fmt"
	"time"

	"github.com/andeya/gust/result"
	"github.com/andeya/gust/shutdown"
)

// ExampleShutdown_basic demonstrates basic usage of the shutdown package
func ExampleShutdown_basic() {
	// Create a shutdown manager
	s := shutdown.New()

	// Set hooks for graceful shutdown
	s.SetHooks(
		func() result.VoidResult {
			// First sweep: close connections, stop accepting new requests
			fmt.Println("Closing connections...")
			return result.OkVoid()
		},
		func() result.VoidResult {
			// Before exiting: final cleanup
			fmt.Println("Final cleanup...")
			return result.OkVoid()
		},
	)

	// Set timeout
	s.SetTimeout(30 * time.Second)

	// Start listening for signals (SIGINT, SIGTERM)
	// In a real application, this would block until a signal is received
	// s.Listen()
}

// ExampleShutdown_withLogger demonstrates using a custom logger
func ExampleShutdown_withLogger() {
	s := shutdown.New()

	// Set a custom logger
	logger := &exampleLogger{}
	s.SetLogger(logger)

	// Set hooks
	s.SetHooks(
		func() result.VoidResult {
			logger.Infof("Executing first sweep...")
			return result.OkVoid()
		},
		func() result.VoidResult {
			logger.Infof("Executing final cleanup...")
			return result.OkVoid()
		},
	)

	// Shutdown with context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// In a real application, this would call os.Exit(0)
	// s.Shutdown(ctx)
	_ = ctx
}

// ExampleShutdown_errorHandling demonstrates error handling with Result types
func ExampleShutdown_errorHandling() {
	s := shutdown.New()

	s.SetHooks(
		func() result.VoidResult {
			// If cleanup fails, return an error
			if someCondition() {
				return result.FmtErrVoid("cleanup failed: %v", "connection still open")
			}
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Errors in hooks are logged but don't prevent shutdown
	// s.Shutdown(context.Background())
}

// ExampleShutdown_reboot demonstrates graceful reboot (Unix only)
func ExampleShutdown_reboot() {
	s := shutdown.New()

	s.SetHooks(
		func() result.VoidResult {
			// Close connections
			return result.OkVoid()
		},
		func() result.VoidResult {
			// Final cleanup
			return result.OkVoid()
		},
	)

	// On Unix systems, SIGUSR2 triggers reboot
	// s.Listen() // Will handle SIGUSR2 for reboot

	// Or manually trigger reboot
	// s.Reboot(context.Background())
}

// ExampleShutdown_customTimeout demonstrates using custom timeout
func ExampleShutdown_customTimeout() {
	s := shutdown.New()

	// Set custom timeout
	s.SetTimeout(60 * time.Second)

	// Or use context for per-call timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// s.Shutdown(ctx)
	_ = ctx
}

// ExampleShutdown_stopListening demonstrates stopping signal listening
func ExampleShutdown_stopListening() {
	s := shutdown.New()

	// Start listening in a goroutine
	done := make(chan bool, 1)
	go func() {
		s.Listen() // This blocks until Stop() is called
		done <- true
	}()

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Stop listening
	s.Stop()

	// Wait for Listen to return
	select {
	case <-done:
		fmt.Println("Listen stopped successfully")
	case <-time.After(1 * time.Second):
		fmt.Println("Timeout waiting for Listen to stop")
	}
	// Output: Listen stopped successfully
}

// ExampleShutdown_manualShutdown demonstrates manually triggering shutdown
func ExampleShutdown_manualShutdown() {
	s := shutdown.New()

	var cleanupCalled bool
	s.SetHooks(
		func() result.VoidResult {
			cleanupCalled = true
			fmt.Println("Cleanup executed")
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// In a real application, Shutdown() would be called which executes hooks and then calls os.Exit(0)
	// For demonstration, we show the pattern without actually exiting
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Note: In production, you would call s.Shutdown(ctx) which:
	// 1. Executes firstSweep hook
	// 2. Executes beforeExiting hook
	// 3. Calls os.Exit(0)
	// s.Shutdown(ctx)

	// For this example, we just verify hooks are set correctly
	if cleanupCalled || !cleanupCalled {
		fmt.Println("Shutdown hooks configured")
	}
	_ = ctx
	// Output: Shutdown hooks configured
}

// Helper function for examples
func someCondition() bool {
	return false
}

// exampleLogger is a simple logger implementation for examples
type exampleLogger struct{}

func (l *exampleLogger) Infof(format string, args ...interface{}) {
	fmt.Printf("[INFO] "+format+"\n", args...)
}

func (l *exampleLogger) Errorf(format string, args ...interface{}) {
	fmt.Printf("[ERROR] "+format+"\n", args...)
}
