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
	o := once.GetValue(true)
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestLazyValue(t *testing.T) {
	assert.Equal(t, gust.None[int](), new(gust.LazyValue[int]).TryGet())
	assert.Equal(t, 0, new(gust.LazyValue[int]).GetValue(false))
	assert.Equal(t, 1, new(gust.LazyValue[int]).Init(1).GetValue(true))
	o := new(one)
	once := new(gust.LazyValue[*one]).InitBySetter(func(ptr **one) error {
		o.Increment()
		*ptr = o
		return nil
	}).Unwrap()
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

func TestLazyValuePanic1(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			assert.Equal(t, "failed", p)
		} else {
			t.Fatalf("should painc")
		}
	}()
	var once = new(gust.LazyValue[struct{}]).InitBySetter(func(*struct{}) error {
		panic("failed")
	}).Unwrap()
	_ = once.TryGet()
	t.Fatalf("unreachable")
}

func TestLazyValuePanic2(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			assert.Equal(t, "LazyValue is not initialized", p)
		} else {
			t.Fatalf("should painc")
		}
	}()
	_ = new(gust.LazyValue[struct{}]).GetValue(true)
	t.Fatalf("unreachable")
}

func BenchmarkLazyValue(b *testing.B) {
	var once = new(gust.LazyValue[struct{}])
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// once.GetValue(false)
			once.TryGet()
		}
	})
}
