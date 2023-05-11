package gust

import (
	"sync"
	"sync/atomic"
)

// NewMutex returns a new *Mutex.
func NewMutex[T any](data T) *Mutex[T] {
	return &Mutex[T]{data: data}
}

// NewRWMutex returns a new *RWMutex.
func NewRWMutex[T any](data T) *RWMutex[T] {
	return &RWMutex[T]{data: data}
}

// Mutex is a better generic-type wrapper for `sync.Mutex` that holds a value.
// A Mutex is a mutual exclusion lock.
// The zero value for a Mutex is an unlocked mutex.
//
// A Mutex must not be copied after first use.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m.
// A successful call to TryLock is equivalent to a call to Lock.
// A failed call to TryLock does not establish any “synchronizes before”
// relation at all.
type Mutex[T any] struct {
	inner sync.Mutex
	data  T
}

// Lock locks m.
// If the lock is already in use, the calling goroutine
// blocks until the mutex is available.
func (m *Mutex[T]) Lock() T {
	m.inner.Lock()
	return m.data
}

// TryLock tries to lock m and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (m *Mutex[T]) TryLock() Option[T] {
	if m.inner.TryLock() {
		return Some(m.data)
	}
	return None[T]()
}

// Unlock unlocks m.
// It is a run-time error if m is not locked on entry to Unlock.
//
// A locked Mutex is not associated with a particular goroutine.
// It is allowed for one goroutine to lock a Mutex and then
// arrange for another goroutine to unlock it.
func (m *Mutex[T]) Unlock(newData ...T) {
	if len(newData) > 0 {
		m.data = newData[0]
	}
	m.inner.Unlock()
}

// LockScope securely read and write the data in the Mutex[T].
func (m *Mutex[T]) LockScope(f func(old T) (new T)) {
	m.inner.Lock()
	defer m.inner.Unlock()
	m.data = f(m.data)
}

// TryLockScope tries to securely read and write the data in the Mutex[T].
func (m *Mutex[T]) TryLockScope(f func(old T) (new T)) {
	if m.inner.TryLock() {
		defer m.inner.Unlock()
		m.data = f(m.data)
	}
}

// RWMutex is a better generic-type wrapper for `sync.RWMutex` that holds a value.
// A RWMutex is a reader/writer mutual exclusion lock.
// The lock can be held by an arbitrary number of readers or a single writer.
// The zero value for a RWMutex is an unlocked mutex.
//
// A RWMutex must not be copied after first use.
//
// If a goroutine holds a RWMutex for reading and another goroutine might
// call Lock, no goroutine should expect to be able to acquire a read lock
// until the initial read lock is released. In particular, this prohibits
// recursive read locking. This is to ensure that the lock eventually becomes
// available; a blocked Lock call excludes new readers from acquiring the
// lock.
//
// In the terminology of the Go memory model,
// the n'th call to Unlock “synchronizes before” the m'th call to Lock
// for any n < m, just as for Mutex.
// For any call to RLock, there exists an n such that
// the n'th call to Unlock “synchronizes before” that call to RLock,
// and the corresponding call to RUnlock “synchronizes before”
// the n+1'th call to Lock.
type RWMutex[T any] struct {
	inner sync.RWMutex
	data  T
}

// Lock locks rw for writing.
// If the lock is already locked for reading or writing,
// Lock blocks until the lock is available.
func (m *RWMutex[T]) Lock() T {
	m.inner.Lock()
	return m.data
}

// TryLock tries to lock rw for writing and reports whether it succeeded.
//
// Note that while correct uses of TryLock do exist, they are rare,
// and use of TryLock is often a sign of a deeper problem
func (m *RWMutex[T]) TryLock() Option[T] {
	if m.inner.TryLock() {
		return Some(m.data)
	}
	return None[T]()
}

// Unlock unlocks rw for writing. It is a run-time error if rw is
// not locked for writing on entry to Unlock.
//
// As with Mutexes, a locked RWMutex is not associated with a particular
// goroutine. One goroutine may RLock (Lock) a RWMutex and then
// arrange for another goroutine to RUnlock (Unlock) it.
func (m *RWMutex[T]) Unlock(newData ...T) {
	if len(newData) > 0 {
		m.data = newData[0]
	}
	m.inner.Unlock()
}

// TryLockScope tries to securely read and write the data in the RWMutex[T].
func (m *RWMutex[T]) TryLockScope(write func(old T) (new T)) {
	if m.inner.TryLock() {
		defer m.inner.Unlock()
		m.data = write(m.data)
	}
}

// LockScope securely read and write the data in the RWMutex[T].
func (m *RWMutex[T]) LockScope(write func(old T) (new T)) {
	m.inner.Lock()
	defer m.inner.Unlock()
	m.data = write(m.data)
}

// Happens-before relationships are indicated to the race detector via:
// - Unlock  -> Lock:  readerSem
// - Unlock  -> RLock: readerSem
// - RUnlock -> Lock:  writerSem
//
// The methods below temporarily disable handling of race synchronization
// events in order to provide the more precise model above to the race
// detector.
//
// For example, atomic.AddInt32 in RLock should not appear to provide
// acquire-release semantics, which would incorrectly synchronize racing
// readers, thus potentially missing races.

// RLock locks rw for reading.
//
// It should not be used for recursive read locking; a blocked Lock
// call excludes new readers from acquiring the lock. See the
// documentation on the RWMutex type.
func (m *RWMutex[T]) RLock() T {
	m.inner.RLock()
	return m.data
}

// TryRLock tries to lock rw for reading and reports whether it succeeded.
//
// Note that while correct uses of TryRLock do exist, they are rare,
// and use of TryRLock is often a sign of a deeper problem
// in a particular use of mutexes.
func (m *RWMutex[T]) TryRLock() Option[T] {
	if m.inner.TryRLock() {
		return Some(m.data)
	}
	return None[T]()
}

// RUnlock undoes a single RLock call;
// it does not affect other simultaneous readers.
// It is a run-time error if rw is not locked for reading
// on entry to RUnlock.
func (m *RWMutex[T]) RUnlock() {
	m.inner.RUnlock()
}

// TryRLockScope tries to securely read the data in the RWMutex[T].
func (m *RWMutex[T]) TryRLockScope(read func(T)) {
	if m.inner.TryRLock() {
		defer m.inner.RUnlock()
		read(m.data)
	}
}

// RLockScope securely read the data in the RWMutex[T].
func (m *RWMutex[T]) RLockScope(read func(T)) {
	m.inner.RLock()
	defer m.inner.RUnlock()
	read(m.data)
}

// TryBest tries to read and do the data in the RWMutex[T] safely,
// swapping the data when readAndDo returns false and then trying to do again.
func (m *RWMutex[T]) TryBest(readAndDo func(T) bool, swapWhenFalse func(old T) (new Option[T])) {
	if readAndDo == nil {
		return
	}
	var ok bool
	m.RLockScope(func(old T) {
		ok = readAndDo(old)
	})
	if ok || swapWhenFalse == nil {
		return
	}
	m.inner.Lock()
	defer m.inner.Unlock()
	if readAndDo(m.data) {
		return
	}
	swapWhenFalse(m.data).Inspect(func(newT T) {
		m.data = newT
		readAndDo(newT)
	})
}

// SyncMap is a better generic-type wrapper for `sync.Map`.
// A SyncMap is like a Go map[interface{}]interface{} but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
//
// The SyncMap type is specialized. Most code should use a plain Go map instead,
// with separate locking or coordination, for better type safety and to make it
// easier to maintain other invariants along with the map content.
//
// The SyncMap type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a SyncMap may significantly reduce lock
// contention compared to a Go map paired with a separate Mutex or RWMutex.
//
// The zero SyncMap is empty and ready for use. A SyncMap must not be copied after first use.
//
// In the terminology of the Go memory model, SyncMap arranges that a write operation
// “synchronizes before” any read operation that observes the effect of the write, where
// read and write operations are defined as follows.
// Load, LoadAndDelete, LoadOrStore are read operations;
// Delete, LoadAndDelete, and Store are write operations;
// and LoadOrStore is a write operation when it returns loaded set to false.
type SyncMap[K any, V any] struct {
	inner sync.Map
}

// Load returns the value stored in the map for a key.
func (m *SyncMap[K, V]) Load(key K) Option[V] {
	return BoolAssertOpt[V](m.inner.Load(key))
}

// Store sets the value for a key.
func (m *SyncMap[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores the given value, and returns None.
func (m *SyncMap[K, V]) LoadOrStore(key K, value V) (existingValue Option[V]) {
	return BoolAssertOpt[V](m.inner.LoadOrStore(key, value))
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
func (m *SyncMap[K, V]) LoadAndDelete(key K) (deletedValue Option[V]) {
	return BoolAssertOpt[V](m.inner.LoadAndDelete(key))
}

// Delete deletes the value for a key.
func (m *SyncMap[K, V]) Delete(key K) {
	m.inner.Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the SyncMap's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *SyncMap[K, V]) Range(f func(key K, value V) bool) {
	m.inner.Range(func(key any, value any) bool {
		k, ok := key.(K)
		if !ok {
			return false
		}
		v, ok := value.(V)
		if !ok {
			return false
		}
		return f(k, v)
	})
}

// AtomicValue is a better generic-type wrapper for `atomic.Value`.
// A AtomicValue provides an atomic load and store of a consistently typed value.
// The zero value for a AtomicValue returns nil from Load.
// Once Store has been called, a AtomicValue must not be copied.
//
// A AtomicValue must not be copied after first use.
type AtomicValue[T any] struct {
	inner atomic.Value
}

// Load returns the value set by the most recent Store.
// It returns None if there has been no call to Store for this AtomicValue.
func (v *AtomicValue[T]) Load() (val Option[T]) {
	return AssertOpt[T](v.inner.Load())
}

// Store sets the value of the AtomicValue to x.
// All calls to Store for a given AtomicValue must use values of the same concrete type.
// Store of an inconsistent type panics, as does Store(nil).
func (v *AtomicValue[T]) Store(val T) {
	v.inner.Store(val)
}

// Swap stores new into AtomicValue and returns the previous value. It returns None if
// the AtomicValue is empty.
//
// All calls to Swap for a given AtomicValue must use values of the same concrete
// type. Swap of an inconsistent type panics, as does Swap(nil).
func (v *AtomicValue[T]) Swap(new T) (old Option[T]) {
	return AssertOpt[T](v.inner.Swap(new))
}

// CompareAndSwap executes the compare-and-swap operation for the AtomicValue.
//
// All calls to CompareAndSwap for a given AtomicValue must use values of the same
// concrete type. CompareAndSwap of an inconsistent type panics, as does
// CompareAndSwap(old, nil).
func (v *AtomicValue[T]) CompareAndSwap(old T, new T) (swapped bool) {
	return v.inner.CompareAndSwap(old, new)
}

// LazyValue a value that can be lazily initialized once and read concurrently.
type LazyValue[T any] struct {
	// done indicates whether the action has been performed.
	// It is first in the struct because it is used in the hot path.
	// The hot path is inlined at every call site.
	// Placing done first allows more compact instructions on some architectures (amd64/386),
	// and fewer instructions (to calculate offset) on other architectures.
	done     uint32
	m        sync.Mutex
	value    Result[T]
	onceInit func(ptr *T) error
}

// SetInitSetter set initialization function.
// NOTE: onceInit can not be nil
func (o *LazyValue[T]) SetInitSetter(onceInit func(ptr *T) error) *LazyValue[T] {
	if o.IsInitialized() {
		return o
	}
	o.m.Lock()
	defer o.m.Unlock()
	o.onceInit = onceInit
	return o
}

// SetInitClosure set initialization function.
// NOTE: onceInit can not be nil
func (o *LazyValue[T]) SetInitClosure(onceInit func() error) *LazyValue[T] {
	return o.SetInitSetter(func(ptr *T) error {
		return onceInit()
	})
}

// SetInitValue set the initialization value.
func (o *LazyValue[T]) SetInitValue(v T) *LazyValue[T] {
	_ = o.SetInitSetter(func(ptr *T) error {
		*ptr = v
		return nil
	})
	return o
}

// IsInitialized determine whether it is initialized.
func (o *LazyValue[T]) IsInitialized() bool {
	return atomic.LoadUint32(&o.done) != 0
}

func (o *LazyValue[T]) markInit() {
	atomic.StoreUint32(&o.done, 1)
}

const ErrLazyValueWithoutInit = "*LazyValue[T]: onceInit function is nil"

// TryGetValue concurrency-safe get the Option[T].
// NOTE: if it is not initialized, return None
func (o *LazyValue[T]) TryGetValue() Result[T] {
	if !o.IsInitialized() {
		o.m.Lock()
		defer o.m.Unlock()
		if o.done == 0 {
			defer o.markInit()
			if o.onceInit == nil {
				o.value = Err[T](ErrLazyValueWithoutInit)
			} else {
				var v T
				err := o.onceInit(&v)
				o.value = Ret[T](v, err)
			}
		}
	}
	return o.value
}
