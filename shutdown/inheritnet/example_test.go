//go:build !windows
// +build !windows

package inheritnet_test

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/andeya/gust/result"
	"github.com/andeya/gust/shutdown"
	"github.com/andeya/gust/shutdown/inheritnet"
)

// Example demonstrates how to use inheritnet with the shutdown package
// for graceful restarts with connection inheritance.
func Example() {
	// Create a shutdown manager
	s := shutdown.New()
	s.SetTimeout(30 * time.Second)

	// Set up hooks for graceful shutdown
	s.SetHooks(
		func() result.VoidResult {
			// First sweep: stop accepting new connections
			// In a real application, you would close the listener here
			log.Println("First sweep: stopping new connections")
			return result.OkVoid()
		},
		func() result.VoidResult {
			// Before exiting: final cleanup
			log.Println("Before exiting: final cleanup")
			return result.OkVoid()
		},
	)

	// Listen on a network address (will inherit if available)
	ln, err := inheritnet.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Set up inheritance before reboot
	// This extracts file descriptors from active listeners and sets them up
	// for inheritance during graceful reboot
	err = inheritnet.SetInherited(s)
	if err != nil {
		log.Fatal(err)
	}

	// Start your HTTP server
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello from server (PID: %d)\n", getPID())
		}),
	}

	// Start server in a goroutine
	go func() {
		if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
			log.Printf("Server error: %v", err)
		}
	}()

	log.Println("Server started on :8080")
	log.Println("Send SIGUSR2 to trigger graceful restart")
	log.Println("Send SIGINT or SIGTERM to trigger graceful shutdown")

	// Start listening for shutdown signals
	// This will block until a signal is received
	s.Listen()
}

// getPID returns the current process ID (for demonstration purposes)
func getPID() int {
	// In a real application, you would use os.Getpid()
	return 0
}

// ExampleSetInherited demonstrates how to set up inheritance for graceful restarts.
func ExampleSetInherited() {
	// Create a shutdown manager
	s := shutdown.New()

	// Listen on a network address
	ln, err := inheritnet.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Set up inheritance before reboot
	// This must be called before calling s.Reboot() to ensure
	// the new process can inherit the active listeners
	err = inheritnet.SetInherited(s)
	if err != nil {
		log.Fatal(err)
	}

	// Now when you call s.Reboot(context.Background()), the new process
	// will inherit the file descriptors from the active listeners
	_ = ln
}

// ExampleListenTCP demonstrates how to use ListenTCP for TCP connections.
func ExampleListenTCP() {
	// Create a shutdown manager
	s := shutdown.New()

	// Listen on a TCP address
	addr, err := inheritnet.ListenTCP("tcp", nil) // nil means use default address
	if err != nil {
		log.Fatal(err)
	}

	// Set up inheritance
	err = inheritnet.SetInherited(s)
	if err != nil {
		log.Fatal(err)
	}

	// Use the listener
	_ = addr
	_ = s
}

// ExampleListenUnix demonstrates how to use ListenUnix for Unix domain sockets.
func ExampleListenUnix() {
	// Create a shutdown manager
	s := shutdown.New()

	// Listen on a Unix domain socket
	addr, err := inheritnet.ListenUnix("unix", nil) // nil means use default address
	if err != nil {
		log.Fatal(err)
	}

	// Set up inheritance
	err = inheritnet.SetInherited(s)
	if err != nil {
		log.Fatal(err)
	}

	// Use the listener
	_ = addr
	_ = s
}

// Example_complete demonstrates a complete graceful restart scenario.
func Example_complete() {
	// Create a shutdown manager
	s := shutdown.New()
	s.SetTimeout(30 * time.Second)

	// Set up hooks
	s.SetHooks(
		func() result.VoidResult {
			// First sweep: stop accepting new connections
			log.Println("Stopping new connections...")
			return result.OkVoid()
		},
		func() result.VoidResult {
			// Before exiting: wait for active connections to close
			log.Println("Waiting for active connections to close...")
			time.Sleep(1 * time.Second) // Simulate waiting
			return result.OkVoid()
		},
	)

	// Listen on network address
	ln, err := inheritnet.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	// Set up inheritance
	err = inheritnet.SetInherited(s)
	if err != nil {
		log.Fatal(err)
	}

	// Start server
	go func() {
		server := &http.Server{Handler: http.DefaultServeMux}
		server.Serve(ln)
	}()

	// Listen for signals in a goroutine
	go s.Listen()

	// In a real application, you would continue with your main logic here
	// For example, you might call s.Reboot(context.Background()) when
	// you want to trigger a graceful restart

	// Keep the program running
	select {}
}

