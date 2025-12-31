// Package coarsetime provides fast, coarse-grained time access.
// It's a faster alternative to time.Now() with configurable precision.
//
// This package is inspired by Linux's CLOCK_REALTIME_COARSE and CLOCK_MONOTONIC_COARSE,
// providing both wall clock time and monotonic time with reduced system call overhead.
//
// # Examples
//
//	```go
//	// Use default instance (100ms precision)
//	now := coarsetime.Now()
//
//	// Create instance with default precision (100ms)
//	ct := coarsetime.New()  // Uses default 100ms precision
//	defer ct.Stop()
//	now := ct.Now()
//
//	// Use custom precision
//	ct2 := coarsetime.New(10 * time.Millisecond)  // 10ms precision
//	defer ct2.Stop()
//	now2 := ct2.Now()
//
//	// Use monotonic time for performance measurement
//	start := coarsetime.Monotonic()
//	// ... do work ...
//	elapsed := coarsetime.Since(start)
//	```
package coarsetime

import (
	"sync/atomic"
	"time"
)

// CoarseTime provides fast, coarse-grained time access.
// It's a faster alternative to time.Now() with configurable precision.
//
// CoarseTime maintains two types of time:
//   - Wall Clock Time: Affected by system clock adjustments, suitable for logging and timestamps
//   - Monotonic Time: Not affected by system clock adjustments, suitable for performance measurement and timeouts
type CoarseTime struct {
	frequency     time.Duration
	wallTime      atomic.Value // *time.Time (wall clock)
	monotonicTime atomic.Value // *time.Time (monotonic, relative to base)
	baseMonotonic time.Time    // Base monotonic time for calculation
	stopCh        chan struct{}
	stopped       atomic.Bool
}

// New creates a new CoarseTime instance with the specified frequency.
// If no frequency is provided, it defaults to 100ms (industry standard for coarse-grained time).
// If frequency <= 0, it also defaults to 100ms.
//
// The frequency determines how often the time is updated. Smaller values provide
// higher precision but use more CPU. Typical values are 10ms, 100ms, or 1s.
//
// Examples:
//   - New() - uses default 100ms precision
//   - New(10 * time.Millisecond) - uses 10ms precision
func New(frequency ...time.Duration) *CoarseTime {
	var freq time.Duration
	if len(frequency) > 0 {
		freq = frequency[0]
	}
	if freq <= 0 {
		freq = 100 * time.Millisecond // default
	}

	ct := &CoarseTime{
		frequency:     freq,
		stopCh:        make(chan struct{}),
		baseMonotonic: time.Now(),
	}

	// Initialize with current time
	now := time.Now()
	truncated := now.Truncate(freq)
	ct.wallTime.Store(&truncated)
	ct.monotonicTime.Store(&truncated)

	// Start update goroutine
	go ct.updateLoop()

	return ct
}

// Now returns the current wall clock time (coarse-grained).
// This is faster than time.Now() but less precise.
//
// The returned time is affected by system clock adjustments.
// For time measurements that should not be affected by clock adjustments,
// use Monotonic() instead.
func (ct *CoarseTime) Now() time.Time {
	if ct.stopped.Load() {
		// If stopped, return current time directly
		return time.Now().Truncate(ct.frequency)
	}
	tp := ct.wallTime.Load().(*time.Time)
	return *tp
}

// Monotonic returns the current monotonic time (coarse-grained).
// This time is not affected by system clock adjustments.
//
// Monotonic time is suitable for:
//   - Performance measurements
//   - Timeout calculations
//   - Duration measurements
//
// Note: Monotonic time should only be used for relative time calculations,
// not for representing absolute wall clock time.
func (ct *CoarseTime) Monotonic() time.Time {
	if ct.stopped.Load() {
		// If stopped, calculate from base
		elapsed := time.Since(ct.baseMonotonic)
		return ct.baseMonotonic.Add(elapsed.Truncate(ct.frequency))
	}
	tp := ct.monotonicTime.Load().(*time.Time)
	return *tp
}

// Floor returns the current time rounded down to the nearest frequency boundary.
// This is equivalent to Now().
func (ct *CoarseTime) Floor() time.Time {
	return ct.Now()
}

// Ceiling returns the current time rounded up to the next frequency boundary.
func (ct *CoarseTime) Ceiling() time.Time {
	now := ct.Now()
	return now.Add(ct.frequency)
}

// Frequency returns the update frequency of this CoarseTime instance.
func (ct *CoarseTime) Frequency() time.Duration {
	return ct.frequency
}

// Since returns the time elapsed since t using monotonic time.
// This is more accurate than time.Since() for coarse-grained measurements
// when system clock adjustments might occur.
func (ct *CoarseTime) Since(t time.Time) time.Duration {
	return ct.Monotonic().Sub(t)
}

// Until returns the duration until t using monotonic time.
func (ct *CoarseTime) Until(t time.Time) time.Duration {
	return t.Sub(ct.Monotonic())
}

// Before reports whether the time instant ct.Now() is before t.
func (ct *CoarseTime) Before(t time.Time) bool {
	return ct.Now().Before(t)
}

// After reports whether the time instant ct.Now() is after t.
func (ct *CoarseTime) After(t time.Time) bool {
	return ct.Now().After(t)
}

// Stop stops the update goroutine and releases resources.
// After calling Stop, the time values will no longer be updated.
// It's safe to call Stop multiple times.
func (ct *CoarseTime) Stop() {
	if ct.stopped.Swap(true) {
		return // already stopped
	}
	close(ct.stopCh)
}

// IsStopped returns whether this CoarseTime instance has been stopped.
func (ct *CoarseTime) IsStopped() bool {
	return ct.stopped.Load()
}

func (ct *CoarseTime) updateLoop() {
	ticker := time.NewTicker(ct.frequency)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			now := time.Now()
			truncated := now.Truncate(ct.frequency)
			ct.wallTime.Store(&truncated)

			// Update monotonic time (relative to base)
			elapsed := now.Sub(ct.baseMonotonic)
			newMonotonic := ct.baseMonotonic.Add(elapsed.Truncate(ct.frequency))
			ct.monotonicTime.Store(&newMonotonic)

		case <-ct.stopCh:
			return
		}
	}
}

// Default instance with 100ms precision (balanced performance and accuracy)
var Default = New(100 * time.Millisecond)

// Common precision presets for convenience
var (
	// Fast10ms provides 10ms precision (higher precision, more CPU usage)
	// Suitable for applications requiring finer-grained time resolution.
	Fast10ms = New(10 * time.Millisecond)

	// Standard100ms provides 100ms precision (balanced, default)
	// Suitable for most applications requiring coarse-grained time.
	Standard100ms = Default

	// Coarse1s provides 1s precision (lowest precision, least CPU usage)
	// Suitable for applications where time precision is not critical.
	Coarse1s = New(1 * time.Second)
)

// Convenience functions using the default instance

// Now returns the current wall clock time using the default instance.
// This is faster than time.Now() but less precise (100ms precision).
func Now() time.Time {
	return Default.Now()
}

// Monotonic returns the current monotonic time using the default instance.
// This time is not affected by system clock adjustments.
func Monotonic() time.Time {
	return Default.Monotonic()
}

// Floor returns the current time rounded down using the default instance.
func Floor() time.Time {
	return Default.Floor()
}

// Ceiling returns the current time rounded up using the default instance.
func Ceiling() time.Time {
	return Default.Ceiling()
}

// Since returns the time elapsed since t using monotonic time from the default instance.
func Since(t time.Time) time.Duration {
	return Default.Since(t)
}

// Until returns the duration until t using monotonic time from the default instance.
func Until(t time.Time) time.Duration {
	return Default.Until(t)
}

// Legacy functions for backward compatibility

// FloorTimeNow returns the current time rounded down to the nearest frequency boundary.
// This is a legacy function for backward compatibility.
// For new code, use Now() or Floor() instead.
func FloorTimeNow() time.Time {
	return Default.Floor()
}

// CeilingTimeNow returns the current time rounded up to the next frequency boundary.
// This is a legacy function for backward compatibility.
// For new code, use Ceiling() instead.
func CeilingTimeNow() time.Time {
	return Default.Ceiling()
}
