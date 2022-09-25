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
func (m *RWMutex[T]) RUnlock(newData ...T) {
	if len(newData) > 0 {
		m.data = newData[0]
	}
	m.inner.RUnlock()
}
