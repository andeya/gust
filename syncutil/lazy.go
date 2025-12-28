package syncutil

import (
	"sync"
	"sync/atomic"

	"github.com/andeya/gust/result"
)

// LazyValue a value that can be lazily initialized once and read concurrently.
type LazyValue[T any] struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	done     uint32
	m        sync.Mutex
	value    result.Result[T]
	onceInit func() result.Result[T]
}

// NewLazyValue new empty LazyValue.
func NewLazyValue[T any]() *LazyValue[T] {
	return new(LazyValue[T])
}

// NewLazyValueWithFunc new LazyValue with initialization function.
// The value will be computed lazily when TryGetValue() is called.
func NewLazyValueWithFunc[T any](onceInit func() result.Result[T]) *LazyValue[T] {
	return new(LazyValue[T]).SetInitFunc(onceInit)
}

// NewLazyValueWithValue new LazyValue with initialization value.
// The value will be computed lazily when TryGetValue() is called.
func NewLazyValueWithValue[T any](v T) *LazyValue[T] {
	return new(LazyValue[T]).SetInitValue(v)
}

// NewLazyValueWithZero new LazyValue with zero.
// The value will be computed lazily when TryGetValue() is called.
func NewLazyValueWithZero[T any]() *LazyValue[T] {
	return new(LazyValue[T]).SetInitZero()
}

// SetInitFunc set initialization function.
// NOTE: onceInit can not be nil
// If the LazyValue already has an initialization function set (even if not initialized yet),
// this function will not override it.
func (o *LazyValue[T]) SetInitFunc(onceInit func() result.Result[T]) *LazyValue[T] {
	if o.IsInitialized() {
		return o
	}
	o.m.Lock()
	defer o.m.Unlock()
	// Don't override if onceInit is already set
	if o.onceInit != nil {
		return o
	}
	o.onceInit = onceInit
	return o
}

// SetInitValue set the initialization value.
func (o *LazyValue[T]) SetInitValue(v T) *LazyValue[T] {
	_ = o.SetInitFunc(func() result.Result[T] {
		return result.Ok(v)
	})
	return o
}

// Zero creates a zero T.
func (*LazyValue[T]) Zero() T {
	var v T
	return v
}

// SetInitZero set the zero value for initialization.
func (o *LazyValue[T]) SetInitZero() *LazyValue[T] {
	return o.SetInitValue(o.Zero())
}

// IsInitialized determine whether it is initialized.
func (o *LazyValue[T]) IsInitialized() bool {
	return atomic.LoadUint32(&o.done) != 0
}

func (o *LazyValue[T]) markInit() {
	atomic.StoreUint32(&o.done, 1)
}

const ErrLazyValueWithoutInit = "*syncutil.LazyValue[T]: onceInit function is nil"

// TryGetValue concurrency-safe get the Result[T].
func (o *LazyValue[T]) TryGetValue() result.Result[T] {
	if !o.IsInitialized() {
		o.m.Lock()
		defer o.m.Unlock()
		if o.done == 0 {
			defer o.markInit()
			if o.onceInit == nil {
				o.value = result.TryErr[T](ErrLazyValueWithoutInit)
			} else {
				o.value = o.onceInit()
			}
		}
	}
	return o.value
}

// GetPtr returns its pointer or nil.
func (o *LazyValue[T]) GetPtr() *T {
	return o.TryGetValue().AsPtr()
}
