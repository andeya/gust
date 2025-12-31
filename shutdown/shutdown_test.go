package shutdown

import (
	"context"
	"sync"
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

	s.logf("test info")
	s.logErrorf("test error")

	assert.Len(t, mockLogger.infos, 1)
	assert.Len(t, mockLogger.errors, 1)
	assert.Contains(t, mockLogger.infos[0], "test info")
	assert.Contains(t, mockLogger.errors[0], "test error")
}

func TestLogger_Nil(t *testing.T) {
	s := New()
	s.SetLogger(nil)

	// Should not panic
	s.logf("test")
	s.logErrorf("test")
}

// mockLogger is a mock implementation of Logger for testing
type mockLogger struct {
	mu     sync.Mutex
	infos  []string
	errors []string
}

func (m *mockLogger) Infof(format string, args ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infos = append(m.infos, format)
}

func (m *mockLogger) Errorf(format string, args ...interface{}) {
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
