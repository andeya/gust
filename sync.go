package gust

import "sync"

// NewMutex returns a new *Mutex.
func NewMutex[T any](data T) *Mutex[T] {
	return &Mutex[T]{data: data}
}

// NewRWMutex returns a new *RWMutex.
func NewRWMutex[T any](data T) *RWMutex[T] {
	return &RWMutex[T]{data: data}
}

// Mutex is a wrapper of `sync.Mutex` that holds a value.
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

// RWMutex is a wrapper of `sync.RWMutex` that holds a value.
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

// Map is like a Go map[interface{}]interface{} but is safe for concurrent use
// by multiple goroutines without additional locking or coordination.
// Loads, stores, and deletes run in amortized constant time.
//
// The Map type is specialized. Most code should use a plain Go map instead,
// with separate locking or coordination, for better type safety and to make it
// easier to maintain other invariants along with the map content.
//
// The Map type is optimized for two common use cases: (1) when the entry for a given
// key is only ever written once but read many times, as in caches that only grow,
// or (2) when multiple goroutines read, write, and overwrite entries for disjoint
// sets of keys. In these two cases, use of a Map may significantly reduce lock
// contention compared to a Go map paired with a separate Mutex or RWMutex.
//
// The zero Map is empty and ready for use. A Map must not be copied after first use.
//
// In the terminology of the Go memory model, Map arranges that a write operation
// “synchronizes before” any read operation that observes the effect of the write, where
// read and write operations are defined as follows.
// Load, LoadAndDelete, LoadOrStore are read operations;
// Delete, LoadAndDelete, and Store are write operations;
// and LoadOrStore is a write operation when it returns loaded set to false.
type Map[K any, V any] struct {
	inner sync.Map
}

// Load returns the value stored in the map for a key, or nil if no
// value is present.
// The ok result indicates whether value was found in the map.
func (m *Map[K, V]) Load(key K) Option[V] {
	return BoolAssertOpt[V](m.inner.Load(key))
}

// Store sets the value for a key.
func (m *Map[K, V]) Store(key K, value V) {
	m.inner.Store(key, value)
}

// LoadOrStore returns the existing value for the key if present.
// Otherwise, it stores the given value, and returns None.
func (m *Map[K, V]) LoadOrStore(key any, value any) (existingValue Option[V]) {
	return BoolAssertOpt[V](m.inner.LoadOrStore(key, value))
}

// LoadAndDelete deletes the value for a key, returning the previous value if any.
func (m *Map[K, V]) LoadAndDelete(key any) (deletedValue Option[V]) {
	return BoolAssertOpt[V](m.inner.LoadAndDelete(key))
}

// Delete deletes the value for a key.
func (m *Map[K, V]) Delete(key any) {
	m.inner.Delete(key)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
//
// Range does not necessarily correspond to any consistent snapshot of the Map's
// contents: no key will be visited more than once, but if the value for any key
// is stored or deleted concurrently (including by f), Range may reflect any
// mapping for that key from any point during the Range call. Range does not
// block other methods on the receiver; even f itself may call any method on m.
//
// Range may be O(N) with the number of elements in the map even if f returns
// false after a constant number of calls.
func (m *Map[K, V]) Range(f func(key K, value V) bool) {
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
