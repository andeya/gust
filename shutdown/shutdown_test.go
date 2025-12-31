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

func init() {
	// Mock exit to prevent actual process termination during tests
	globalDeps.exit = func(code int) {}
}

func TestNew(t *testing.T) {
	s := New()
	assert.NotNil(t, s)
	assert.Equal(t, MinShutdownTimeout, s.Timeout())
	assert.False(t, s.IsListening())
}

func TestSetTimeout(t *testing.T) {
	s := New()

	// Valid timeout
	s.SetTimeout(30 * time.Second)
	assert.Equal(t, 30*time.Second, s.Timeout())

	// Timeout less than minimum
	s.SetTimeout(5 * time.Second)
	assert.Equal(t, MinShutdownTimeout, s.Timeout())

	// Negative timeout (effectively infinite)
	s.SetTimeout(-1)
	assert.True(t, s.Timeout() > time.Hour*24*365)

	// Zero timeout
	s.SetTimeout(0)
	assert.Equal(t, MinShutdownTimeout, s.Timeout())
}

func TestSetHooks(t *testing.T) {
	s := New()

	// Test with valid hooks
	var firstCalled, beforeCalled bool
	s.SetHooks(
		func() result.VoidResult { firstCalled = true; return result.OkVoid() },
		func() result.VoidResult { beforeCalled = true; return result.OkVoid() },
	)
	s.getFirstSweep()()
	s.getBeforeExiting()()
	assert.True(t, firstCalled)
	assert.True(t, beforeCalled)

	// Test with nil hooks (should use no-op)
	s.SetHooks(nil, nil)
	assert.True(t, s.getFirstSweep()().IsOk())
	assert.True(t, s.getBeforeExiting()().IsOk())
}

func TestLogger(t *testing.T) {
	s := New()

	// Test nil logger (should not panic)
	s.SetLogger(nil)
	s.logf(context.Background(), "test")
	s.logErrorf(context.Background(), "test")

	// Test with mock logger
	ml := &mockLogger{}
	s.SetLogger(ml)
	assert.Equal(t, ml, s.Logger())

	s.logf(context.Background(), "info")
	s.logErrorf(context.Background(), "error")
	assert.Len(t, ml.infos, 1)
	assert.Len(t, ml.errors, 1)
}

func TestExecuteShutdown(t *testing.T) {
	s := New()

	// Success case
	var called bool
	s.SetHooks(
		func() result.VoidResult { called = true; return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)
	res := s.executeShutdown(context.Background())
	assert.True(t, res.IsOk())
	assert.True(t, called)

	// FirstSweep error
	s.SetHooks(
		func() result.VoidResult { return result.FmtErrVoid("first error") },
		func() result.VoidResult { return result.OkVoid() },
	)
	res = s.executeShutdown(context.Background())
	assert.True(t, res.IsErr())

	// BeforeExiting error
	s.SetHooks(
		func() result.VoidResult { return result.OkVoid() },
		func() result.VoidResult { return result.FmtErrVoid("before error") },
	)
	res = s.executeShutdown(context.Background())
	assert.True(t, res.IsErr())
}

func TestShutdownInternal(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult { called = true; return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	// Test with nil context (creates timeout context internally)
	//nolint:staticcheck // intentionally pass nil to test nil context path
	s.shutdownInternal(nil)
	assert.True(t, called)

	// Test with provided context
	called = false
	s.shutdownInternal(context.Background())
	assert.True(t, called)

	// Test timeout path
	s.SetHooks(
		func() result.VoidResult { time.Sleep(100 * time.Millisecond); return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	s.shutdownInternal(ctx) // Should handle timeout gracefully
}

func TestShutdown(t *testing.T) {
	s := New()
	s.SetTimeout(50 * time.Millisecond)

	var called bool
	s.SetHooks(
		func() result.VoidResult { called = true; return result.OkVoid() },
		func() result.VoidResult { return result.OkVoid() },
	)

	s.Shutdown(context.Background())
	assert.True(t, called)
}

func TestStop(t *testing.T) {
	s := New()

	// Stop when not listening
	s.Stop()
	assert.False(t, s.IsListening())

	// Stop when already stopped
	atomic.StoreInt32(&s.stopped, 1)
	s.Stop()

	// Stop when listening
	atomic.StoreInt32(&s.listening, 1)
	atomic.StoreInt32(&s.stopped, 0)
	s.stopCh = make(chan struct{})
	s.Stop()
	assert.False(t, s.IsListening())
	assert.Equal(t, int32(1), atomic.LoadInt32(&s.stopped))

	// Double-check path (listening but stopped)
	atomic.StoreInt32(&s.listening, 1)
	atomic.StoreInt32(&s.stopped, 1)
	s.stopCh = make(chan struct{})
	s.Stop()
}

func TestConcurrentAccess(t *testing.T) {
	s := New()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			s.SetTimeout(time.Duration(i) * time.Millisecond)
			_ = s.Timeout()
			s.SetHooks(nil, nil)
			_ = s.getFirstSweep()()
			_ = s.getBeforeExiting()()
			_ = s.IsListening()
		}(i)
	}
	wg.Wait()
}

// mockLogger is a mock implementation of Logger for testing
type mockLogger struct {
	mu     sync.Mutex
	infos  []string
	errors []string
}

func (m *mockLogger) Infof(_ context.Context, format string, _ ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.infos = append(m.infos, format)
}

func (m *mockLogger) Errorf(_ context.Context, format string, _ ...interface{}) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errors = append(m.errors, format)
}
