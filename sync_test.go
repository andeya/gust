package gust_test

import (
	"testing"

	"github.com/andeya/gust"
	"github.com/stretchr/testify/assert"
)

func TestMutex(t *testing.T) {
	var m = gust.NewMutex(1)
	assert.Equal(t, 1, m.Lock())
	m.Unlock()
	assert.Equal(t, 1, m.Lock())
	m.Unlock(2)
	assert.Equal(t, 2, m.Lock())
	m.Unlock()
}

func TestSyncMap(t *testing.T) {
	var m gust.SyncMap[string, int]
	assert.Equal(t, gust.None[int](), m.Load("a"))
	m.Store("a", 1)
	assert.Equal(t, gust.Some(1), m.Load("a"))
	m.Delete("a")
	assert.Equal(t, gust.None[int](), m.Load("a"))
}

func TestAtomicValue(t *testing.T) {
	var m gust.AtomicValue[int]
	assert.Equal(t, gust.None[int](), m.Load())
	m.Store(1)
	assert.Equal(t, gust.Some(1), m.Load())
	assert.Equal(t, gust.Some(1), m.Swap(2))
	assert.Equal(t, gust.Some(2), m.Load())
	assert.False(t, m.CompareAndSwap(1, 3))
	assert.True(t, m.CompareAndSwap(2, 3))
}

type one int

func (o *one) Increment() {
	*o++
}

func runLazyValue(t *testing.T, once *gust.LazyValue[*one], c chan bool) {
	o := once.Get().Unwrap()
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestLazyValue(t *testing.T) {
	assert.EqualError(t, new(gust.LazyValue[int]).Get().Err(), "LazyValue is not initialized")
	assert.Equal(t, 0, new(gust.LazyValue[int]).UncheckGet())
	assert.Equal(t, 1, new(gust.LazyValue[int]).InitOnce(1).Get().Unwrap())
	o := new(one)
	once := new(gust.LazyValue[*one]).InitOnceBy(func() gust.Result[*one] {
		o.Increment()
		return gust.Ok(o)
	})
	c := make(chan bool)
	const N = 10
	for i := 0; i < N; i++ {
		go runLazyValue(t, once, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 1 {
		t.Errorf("once failed outside run: %d is not 1", *o)
	}
}

func TestLazyValuePanic(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			assert.Equal(t, "failed", p)
		} else {
			t.Fatalf("should painc")
		}
	}()
	var once = new(gust.LazyValue[struct{}]).InitOnceBy(func() gust.Result[struct{}] {
		panic("failed")
	})
	_ = once.Get()
	t.Fatalf("unreachable")
}

func BenchmarkLazyValue(b *testing.B) {
	var once = new(gust.LazyValue[struct{}]).InitZeroOnce()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			once.UncheckGet()
		}
	})
}
