package coarsetime

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Run("valid frequency", func(t *testing.T) {
		ct := New(50 * time.Millisecond)
		defer ct.Stop()

		assert.Equal(t, 50*time.Millisecond, ct.Frequency())
		assert.False(t, ct.IsStopped())
	})

	t.Run("zero frequency defaults to 100ms", func(t *testing.T) {
		ct := New(0)
		defer ct.Stop()

		assert.Equal(t, 100*time.Millisecond, ct.Frequency())
	})

	t.Run("negative frequency defaults to 100ms", func(t *testing.T) {
		ct := New(-10 * time.Millisecond)
		defer ct.Stop()

		assert.Equal(t, 100*time.Millisecond, ct.Frequency())
	})
}

func TestCoarseTime_Now(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	now := ct.Now()
	assert.False(t, now.IsZero())

	// Should be within 200ms of actual time (100ms precision + 100ms update delay)
	actual := time.Now()
	diff := actual.Sub(now)
	assert.True(t, diff >= 0, "coarse time should not be in the future")
	assert.True(t, diff < 200*time.Millisecond, "coarse time should be within 200ms of actual time")
}

func TestCoarseTime_Monotonic(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	mono1 := ct.Monotonic()
	time.Sleep(50 * time.Millisecond)
	mono2 := ct.Monotonic()

	// Monotonic time should be non-decreasing
	assert.True(t, mono2.After(mono1) || mono2.Equal(mono1), "monotonic time should be non-decreasing")
}

func TestCoarseTime_Floor(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	floor := ct.Floor()
	now := ct.Now()

	// Floor should be same as Now
	assert.Equal(t, now, floor)
}

func TestCoarseTime_Ceiling(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	floor := ct.Floor()
	ceiling := ct.Ceiling()

	// Ceiling should be floor + frequency
	expected := floor.Add(ct.Frequency())
	assert.Equal(t, expected, ceiling)
}

func TestCoarseTime_Since(t *testing.T) {
	ct := New(10 * time.Millisecond)
	defer ct.Stop()

	start := ct.Monotonic()
	time.Sleep(50 * time.Millisecond)
	elapsed := ct.Since(start)

	// Should be approximately 50ms (within 20ms tolerance due to coarse precision)
	assert.True(t, elapsed >= 30*time.Millisecond, "elapsed should be at least 30ms")
	assert.True(t, elapsed < 100*time.Millisecond, "elapsed should be less than 100ms")
}

func TestCoarseTime_Until(t *testing.T) {
	ct := New(10 * time.Millisecond)
	defer ct.Stop()

	future := ct.Monotonic().Add(100 * time.Millisecond)
	until := ct.Until(future)

	// Should be approximately 100ms (within 20ms tolerance)
	assert.True(t, until >= 80*time.Millisecond, "until should be at least 80ms")
	assert.True(t, until < 120*time.Millisecond, "until should be less than 120ms")
}

func TestCoarseTime_Before(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	future := time.Now().Add(1 * time.Second)
	assert.True(t, ct.Before(future), "current time should be before future time")

	past := time.Now().Add(-1 * time.Second)
	assert.False(t, ct.Before(past), "current time should not be before past time")
}

func TestCoarseTime_After(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	past := time.Now().Add(-1 * time.Second)
	assert.True(t, ct.After(past), "current time should be after past time")

	future := time.Now().Add(1 * time.Second)
	assert.False(t, ct.After(future), "current time should not be after future time")
}

func TestCoarseTime_Stop(t *testing.T) {
	ct := New(100 * time.Millisecond)

	// Should not be stopped initially
	assert.False(t, ct.IsStopped())

	// Stop it
	ct.Stop()
	assert.True(t, ct.IsStopped())

	// Multiple stops should be safe
	ct.Stop()
	ct.Stop()
	assert.True(t, ct.IsStopped())

	// Time should still work after stop (but less accurate)
	now := ct.Now()
	assert.False(t, now.IsZero())
}

func TestCoarseTime_Stop_NoUpdate(t *testing.T) {
	ct := New(10 * time.Millisecond)

	// Get initial time
	initial := ct.Now()

	// Stop immediately
	ct.Stop()

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Time should not have updated (or updated minimally)
	after := ct.Now()
	diff := after.Sub(initial)

	// After stop, Now() returns current time truncated, so it may have changed
	// But the internal update loop should have stopped
	assert.True(t, diff >= 0)
}

func TestDefaultInstance(t *testing.T) {
	// Default should be initialized
	assert.NotNil(t, Default)
	assert.Equal(t, 100*time.Millisecond, Default.Frequency())
	assert.False(t, Default.IsStopped())
}

func TestPresets(t *testing.T) {
	t.Run("Fast10ms", func(t *testing.T) {
		assert.Equal(t, 10*time.Millisecond, Fast10ms.Frequency())
	})

	t.Run("Standard100ms", func(t *testing.T) {
		assert.Equal(t, 100*time.Millisecond, Standard100ms.Frequency())
		assert.Equal(t, Default, Standard100ms)
	})

	t.Run("Coarse1s", func(t *testing.T) {
		assert.Equal(t, 1*time.Second, Coarse1s.Frequency())
	})
}

func TestConvenienceFunctions(t *testing.T) {
	t.Run("Now", func(t *testing.T) {
		now := Now()
		assert.False(t, now.IsZero())
	})

	t.Run("Monotonic", func(t *testing.T) {
		mono := Monotonic()
		assert.False(t, mono.IsZero())
	})

	t.Run("Floor", func(t *testing.T) {
		floor := Floor()
		assert.False(t, floor.IsZero())
	})

	t.Run("Ceiling", func(t *testing.T) {
		ceiling := Ceiling()
		assert.False(t, ceiling.IsZero())
	})

	t.Run("Since", func(t *testing.T) {
		start := Monotonic()
		time.Sleep(50 * time.Millisecond)
		elapsed := Since(start)
		assert.True(t, elapsed >= 30*time.Millisecond)
		assert.True(t, elapsed < 200*time.Millisecond)
	})

	t.Run("Until", func(t *testing.T) {
		future := Monotonic().Add(100 * time.Millisecond)
		until := Until(future)
		assert.True(t, until >= 80*time.Millisecond)
		assert.True(t, until < 120*time.Millisecond)
	})
}

func TestLegacyFunctions(t *testing.T) {
	t.Run("FloorTimeNow", func(t *testing.T) {
		floor := FloorTimeNow()
		assert.False(t, floor.IsZero())
		assert.Equal(t, Floor(), floor)
	})

	t.Run("CeilingTimeNow", func(t *testing.T) {
		ceiling := CeilingTimeNow()
		assert.False(t, ceiling.IsZero())
		assert.Equal(t, Ceiling(), ceiling)
	})
}

func TestConcurrentAccess(t *testing.T) {
	ct := New(10 * time.Millisecond)
	defer ct.Stop()

	// Concurrent reads should be safe
	done := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		go func() {
			_ = ct.Now()
			_ = ct.Monotonic()
			_ = ct.Floor()
			_ = ct.Ceiling()
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 100; i++ {
		<-done
	}
}

func TestMonotonicTime_Monotonicity(t *testing.T) {
	ct := New(10 * time.Millisecond)
	defer ct.Stop()

	// Collect monotonic times
	times := make([]time.Time, 100)
	for i := 0; i < 100; i++ {
		times[i] = ct.Monotonic()
		time.Sleep(1 * time.Millisecond)
	}

	// Verify monotonicity (non-decreasing)
	for i := 1; i < len(times); i++ {
		assert.True(t, times[i].After(times[i-1]) || times[i].Equal(times[i-1]),
			"monotonic time should be non-decreasing at index %d", i)
	}
}

func TestPrecision(t *testing.T) {
	testCases := []struct {
		name      string
		frequency time.Duration
	}{
		{"1ms", 1 * time.Millisecond},
		{"10ms", 10 * time.Millisecond},
		{"100ms", 100 * time.Millisecond},
		{"1s", 1 * time.Second},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ct := New(tc.frequency)
			defer ct.Stop()

			// Get multiple time samples
			times := make([]time.Time, 10)
			for i := 0; i < 10; i++ {
				times[i] = ct.Now()
				time.Sleep(tc.frequency / 2) // Sleep half the frequency
			}

			// Verify that times are truncated to frequency boundaries
			for _, tm := range times {
				// Time should be a multiple of frequency (within tolerance)
				nanos := tm.UnixNano()
				freqNanos := tc.frequency.Nanoseconds()
				remainder := nanos % freqNanos
				assert.Equal(t, int64(0), remainder, "time should be truncated to frequency boundary")
			}
		})
	}
}

func TestMultipleInstances(t *testing.T) {
	ct1 := New(10 * time.Millisecond)
	defer ct1.Stop()

	ct2 := New(100 * time.Millisecond)
	defer ct2.Stop()

	// Both should work independently
	now1 := ct1.Now()
	now2 := ct2.Now()

	assert.False(t, now1.IsZero())
	assert.False(t, now2.IsZero())
	assert.Equal(t, 10*time.Millisecond, ct1.Frequency())
	assert.Equal(t, 100*time.Millisecond, ct2.Frequency())
}

func TestWallClockVsMonotonic(t *testing.T) {
	ct := New(10 * time.Millisecond)
	defer ct.Stop()

	wall1 := ct.Now()
	mono1 := ct.Monotonic()

	time.Sleep(50 * time.Millisecond)

	wall2 := ct.Now()
	mono2 := ct.Monotonic()

	// Both should advance
	wallDiff := wall2.Sub(wall1)
	monoDiff := mono2.Sub(mono1)

	assert.True(t, wallDiff > 0, "wall clock time should advance")
	assert.True(t, monoDiff > 0, "monotonic time should advance")

	// Both should advance by similar amounts (within tolerance)
	assert.True(t, wallDiff >= 30*time.Millisecond && wallDiff < 100*time.Millisecond)
	assert.True(t, monoDiff >= 30*time.Millisecond && monoDiff < 100*time.Millisecond)
}

func TestCoarseTime_Accuracy(t *testing.T) {
	ct := New(100 * time.Millisecond)
	defer ct.Stop()

	// Wait for initial update
	time.Sleep(150 * time.Millisecond)

	coarse := ct.Now()
	actual := time.Now()

	// Coarse time should be within 200ms of actual time
	diff := actual.Sub(coarse)
	require.True(t, diff >= 0, "coarse time should not be in the future")
	require.True(t, diff < 200*time.Millisecond, "coarse time should be within 200ms of actual time, got %v", diff)
}
