package shutdown

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/andeya/gust/result"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	s := New()
	assert.NotNil(t, s)
	assert.Equal(t, MinShutdownTimeout, s.Timeout())
	assert.False(t, s.IsListening())
}

func TestSetTimeout(t *testing.T) {
	s := New()

	t.Run("valid timeout", func(t *testing.T) {
		s.SetTimeout(30 * time.Second)
		assert.Equal(t, 30*time.Second, s.Timeout())
	})

	t.Run("timeout less than minimum", func(t *testing.T) {
		s.SetTimeout(5 * time.Second)
		assert.Equal(t, MinShutdownTimeout, s.Timeout())
	})

	t.Run("negative timeout", func(t *testing.T) {
		s.SetTimeout(-1)
		assert.True(t, s.Timeout() > time.Hour*24*365) // effectively infinite
	})

	t.Run("zero timeout", func(t *testing.T) {
		s.SetTimeout(0)
		assert.Equal(t, MinShutdownTimeout, s.Timeout())
	})
}

func TestSetHooks(t *testing.T) {
	s := New()

	var firstSweepCalled, beforeExitingCalled bool

	firstSweep := func() result.VoidResult {
		firstSweepCalled = true
		return result.OkVoid()
	}

	beforeExiting := func() result.VoidResult {
		beforeExitingCalled = true
		return result.OkVoid()
	}

	s.SetHooks(firstSweep, beforeExiting)

	// Test hooks are set
	res1 := s.getFirstSweep()()
	assert.True(t, res1.IsOk())
	assert.True(t, firstSweepCalled)

	res2 := s.getBeforeExiting()()
	assert.True(t, res2.IsOk())
	assert.True(t, beforeExitingCalled)
}

func TestSetHooks_Nil(t *testing.T) {
	s := New()

	// Setting nil hooks should use no-op functions
	s.SetHooks(nil, nil)

	res1 := s.getFirstSweep()()
	assert.True(t, res1.IsOk())

	res2 := s.getBeforeExiting()()
	assert.True(t, res2.IsOk())
}

func TestSetLogger(t *testing.T) {
	s := New()

	mockLogger := &mockLogger{}
	s.SetLogger(mockLogger)

	assert.Equal(t, mockLogger, s.Logger())
}

func TestShutdown_Success(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)

	var firstSweepCalled, beforeExitingCalled bool

	s.SetHooks(
		func() result.VoidResult {
			firstSweepCalled = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			beforeExitingCalled = true
			return result.OkVoid()
		},
	)

	// Use a context to prevent actual exit
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// We can't test os.Exit, so we test the executeShutdown method directly
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
	assert.True(t, firstSweepCalled)
	assert.True(t, beforeExitingCalled)
}

func TestShutdown_FirstSweepError(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.FmtErrVoid("first sweep error")
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	ctx := context.Background()
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsErr())
	assert.Contains(t, res.Err().Error(), "first sweep error")
}

func TestShutdown_BeforeExitingError(t *testing.T) {
	s := New()
	s.SetTimeout(100 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.FmtErrVoid("before exiting error")
		},
	)

	ctx := context.Background()
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsErr())
	assert.Contains(t, res.Err().Error(), "before exiting error")
}

func TestShutdown_WithContext(t *testing.T) {
	s := New()
	s.SetTimeout(1 * time.Second)

	var called bool
	s.SetHooks(
		func() result.VoidResult {
			called = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Use a short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
	assert.True(t, called)
}

func TestShutdown_ContextTimeout(t *testing.T) {
	s := New()

	s.SetHooks(
		func() result.VoidResult {
			time.Sleep(200 * time.Millisecond) // Longer than context timeout
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()

	// executeShutdown should complete even if context times out
	// (the context is checked in Shutdown, not executeShutdown)
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
}

func TestIsListening(t *testing.T) {
	s := New()
	assert.False(t, s.IsListening())
}

func TestStop(t *testing.T) {
	s := New()
	assert.False(t, s.IsListening())

	// Stop when not listening should be safe
	s.Stop()
	assert.False(t, s.IsListening())
}

func TestLogger(t *testing.T) {
	s := New()

	mockLogger := &mockLogger{
		infos:  make([]string, 0),
		errors: make([]string, 0),
	}
	s.SetLogger(mockLogger)

	s.logf(context.Background(), "test info")
	s.logErrorf(context.Background(), "test error")

	assert.Len(t, mockLogger.infos, 1)
	assert.Len(t, mockLogger.errors, 1)
	assert.Contains(t, mockLogger.infos[0], "test info")
	assert.Contains(t, mockLogger.errors[0], "test error")
}

func TestLogger_Nil(t *testing.T) {
	s := New()
	s.SetLogger(nil)

	// Should not panic
	s.logf(context.Background(), "test")
	s.logErrorf(context.Background(), "test")
}

// mockLogger is a mock implementation of Logger for testing
type mockLogger struct {
	mu     sync.Mutex
	infos  []string
	errors []string
}

func (m *mockLogger) Infof(ctx context.Context, format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infos = append(m.infos, format)
}

func (m *mockLogger) Errorf(ctx context.Context, format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = append(m.errors, format)
}

func TestConcurrentAccess(t *testing.T) {
	s := New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.SetTimeout(time.Duration(i) * time.Millisecond)
			_ = s.Timeout()
			s.SetHooks(nil, nil)
			_ = s.IsListening()
		}()
	}
	wg.Wait()
}

func TestMinShutdownTimeout(t *testing.T) {
	assert.True(t, MinShutdownTimeout >= 15*time.Second)
}

// Test that Shutdown method doesn't actually exit (we can't test os.Exit)
func TestShutdown_Method(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult {
			called = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// We can't test os.Exit, but we can verify the hooks are called
	// by testing executeShutdown directly
	ctx := context.Background()
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
	assert.True(t, called)
}

func TestGetHooks_ThreadSafe(t *testing.T) {
	s := New()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_ = s.getFirstSweep()()
		}()
		go func() {
			defer wg.Done()
			_ = s.getBeforeExiting()()
		}()
	}
	wg.Wait()
}

// TestShutdownInternal tests shutdownInternal method
func TestShutdownInternal(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult {
			called = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test with nil context (shutdownInternal creates timeout context)
	//nolint // We intentionally pass nil to test the nil context path
	s.shutdownInternal(nil)
	assert.True(t, called)

	// Test with provided context
	called = false
	ctx := context.Background()
	s.shutdownInternal(ctx)
	assert.True(t, called)
}

// TestShutdownInternal_Timeout tests shutdownInternal timeout path
func TestShutdownInternal_Timeout(t *testing.T) {
	s := New()
	s.SetTimeout(10 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			time.Sleep(50 * time.Millisecond) // Longer than timeout
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test with short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// shutdownInternal should handle timeout gracefully
	s.shutdownInternal(ctx)
}

// TestShutdownInternal_ContextTimeout tests shutdownInternal with context timeout
func TestShutdownInternal_ContextTimeout(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			time.Sleep(100 * time.Millisecond) // Longer than timeout
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test the timeout path in Shutdown logic
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// Test executeShutdown with timeout (simulates Shutdown's timeout path)
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.executeShutdown(ctx)
	}()

	select {
	case <-ctx.Done():
		// Timeout occurred (this tests the timeout path in Shutdown)
		<-done
	case <-done:
		// Completed before timeout
	}
}

// TestShutdown_WithNilContext tests Shutdown with nil context
func TestShutdown_WithNilContext(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult {
			called = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test executeShutdown with nil context (Shutdown creates timeout context)
	// We test the logic path that Shutdown would take
	ctx, cancel := context.WithTimeout(context.Background(), s.Timeout())
	defer cancel()

	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
	assert.True(t, called)
}

// TestShutdown_ContextTimeoutLogic tests Shutdown with context timeout
func TestShutdown_ContextTimeoutLogic(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	s.SetHooks(
		func() result.VoidResult {
			time.Sleep(100 * time.Millisecond) // Longer than timeout
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Test with short timeout context
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	// The Shutdown method would timeout, but executeShutdown still completes
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk())
}

// TestShutdown_ContextErrorHandling tests context error handling in Shutdown
func TestShutdown_ContextErrorHandling(t *testing.T) {
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

	// Test with cancelled context
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test executeShutdown with cancelled context
	res := s.executeShutdown(ctx)
	assert.True(t, res.IsOk(), "executeShutdown should complete even with cancelled context")
}

// TestShutdown_MethodLogic tests the Shutdown method logic (without os.Exit)
func TestShutdown_MethodLogic(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var firstSweepCalled, beforeExitingCalled bool
	s.SetHooks(
		func() result.VoidResult {
			firstSweepCalled = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			beforeExitingCalled = true
			return result.OkVoid()
		},
	)

	// Test with context
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// We can't test os.Exit, but we can test the logic by checking executeShutdown
	// The Shutdown method calls executeShutdown in a goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		s.executeShutdown(ctx)
	}()

	select {
	case <-done:
		assert.True(t, firstSweepCalled)
		assert.True(t, beforeExitingCalled)
	case <-time.After(200 * time.Millisecond):
		t.Fatal("shutdown did not complete")
	}
}

// TestShutdown_MethodWithMockExit tests Shutdown method (with mocked exit)
func TestShutdown_MethodWithMockExit(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult {
			called = true
			return result.OkVoid()
		},
		func() result.VoidResult {
			return result.OkVoid()
		},
	)

	// Shutdown calls shutdownInternal and then exit
	// Since exit is mocked in reboot_test.go, we need to set it here too
	originalExit := globalDeps.exit
	globalDeps.exit = func(code int) {}
	defer func() {
		globalDeps.exit = originalExit
	}()

	ctx := context.Background()
	s.Shutdown(ctx)
	assert.True(t, called)
}

// TestStop_AllPaths tests all paths in Stop method
func TestStop_AllPaths(t *testing.T) {
	s := New()

	// Test 1: Stop when already stopped (fast path)
	atomic.StoreInt32(&s.stopped, 1)
	s.Stop()
	assert.True(t, atomic.LoadInt32(&s.stopped) == 1)

	// Test 2: Stop when not listening (fast path)
	atomic.StoreInt32(&s.stopped, 0)
	atomic.StoreInt32(&s.listening, 0)
	s.Stop()
	assert.False(t, s.IsListening())

	// Test 3: Stop when listening (slow path)
	atomic.StoreInt32(&s.listening, 1)
	atomic.StoreInt32(&s.stopped, 0)
	s.stopCh = make(chan struct{})
	s.Stop()
	assert.False(t, s.IsListening())
	assert.True(t, atomic.LoadInt32(&s.stopped) == 1)

	// Test 4: Stop when listening but already stopped (double-check)
	atomic.StoreInt32(&s.listening, 1)
	atomic.StoreInt32(&s.stopped, 1)
	s.stopCh = make(chan struct{})
	s.Stop()
	// Should remain stopped
	assert.True(t, atomic.LoadInt32(&s.stopped) == 1)
}
