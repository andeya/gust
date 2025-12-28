package syncutil

import (
	"sync"

	"github.com/andeya/gust"
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
// the n'th call to Unlock "synchronizes before" the m'th call to Lock
// for any n < m.
// A successful call to TryLock is equivalent to a call to Lock.
// A failed call to TryLock does not establish any "synchronizes before"
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
func (m *Mutex[T]) TryLock() gust.Option[T] {
	if m.inner.TryLock() {
		return gust.Some(m.data)
	}
	return gust.None[T]()
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
// the n'th call to Unlock "synchronizes before" the m'th call to Lock
// for any n < m, just as for Mutex.
// For any call to RLock, there exists an n such that
// the n'th call to Unlock "synchronizes before" that call to RLock,
// and the corresponding call to RUnlock "synchronizes before"
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
func (m *RWMutex[T]) TryLock() gust.Option[T] {
	if m.inner.TryLock() {
		return gust.Some(m.data)
	}
	return gust.None[T]()
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
func (m *RWMutex[T]) TryRLock() gust.Option[T] {
	if m.inner.TryRLock() {
		return gust.Some(m.data)
	}
	return gust.None[T]()
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
func (m *RWMutex[T]) TryBest(readAndDo func(T) bool, swapWhenFalse func(old T) (new gust.Option[T])) {
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
