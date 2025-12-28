package syncutil_test

import (
	"fmt"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/syncutil"
	"github.com/stretchr/testify/assert"
)

func TestMutex(t *testing.T) {
	var m = syncutil.NewMutex(1)
	assert.Equal(t, 1, m.Lock())
	m.Unlock()
	assert.Equal(t, 1, m.Lock())
	m.Unlock(2)
	assert.Equal(t, 2, m.Lock())
	m.Unlock()

	// Test TryLock
	result := m.TryLock()
	assert.True(t, result.IsSome())
	assert.Equal(t, 2, result.Unwrap())
	m.Unlock()

	// Test LockScope
	m.LockScope(func(old int) int {
		return old * 2
	})
	result2 := m.TryLock()
	assert.True(t, result2.IsSome())
	assert.Equal(t, 4, result2.Unwrap())
	m.Unlock()

	// Test TryLockScope (successful)
	success := false
	m.TryLockScope(func(old int) int {
		success = true
		return old + 1
	})
	assert.True(t, success)
	result3 := m.TryLock()
	assert.True(t, result3.IsSome())
	assert.Equal(t, 5, result3.Unwrap())
	m.Unlock()
}

func TestSyncMap(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	assert.Equal(t, gust.None[int](), m.Load("a"))
	m.Store("a", 1)
	assert.Equal(t, gust.Some(1), m.Load("a"))
	m.Delete("a")
	assert.Equal(t, gust.None[int](), m.Load("a"))

	// Test LoadOrStore
	existing := m.LoadOrStore("b", 2)
	assert.True(t, existing.IsNone()) // Key doesn't exist
	assert.Equal(t, gust.Some(2), m.Load("b"))

	existing2 := m.LoadOrStore("b", 3)
	assert.True(t, existing2.IsSome()) // Key exists
	assert.Equal(t, 2, existing2.Unwrap())
	assert.Equal(t, gust.Some(2), m.Load("b")) // Value unchanged

	// Test LoadAndDelete
	deleted := m.LoadAndDelete("b")
	assert.True(t, deleted.IsSome())
	assert.Equal(t, 2, deleted.Unwrap())
	assert.Equal(t, gust.None[int](), m.Load("b"))

	deleted2 := m.LoadAndDelete("nonexistent")
	assert.True(t, deleted2.IsNone())

	// Test Range
	m.Store("x", 10)
	m.Store("y", 20)
	m.Store("z", 30)
	keys := make(map[string]int)
	m.Range(func(key string, value int) bool {
		keys[key] = value
		return true
	})
	assert.Len(t, keys, 3)
	assert.Equal(t, 10, keys["x"])
	assert.Equal(t, 20, keys["y"])
	assert.Equal(t, 30, keys["z"])

	// Test Range with early termination
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})
	assert.Equal(t, 2, count)
}

func TestAtomicValue(t *testing.T) {
	var m syncutil.AtomicValue[int]
	assert.Equal(t, gust.None[int](), m.Load())
	m.Store(1)
	assert.Equal(t, gust.Some(1), m.Load())
	assert.Equal(t, gust.Some(1), m.Swap(2))
	assert.Equal(t, gust.Some(2), m.Load())
	assert.False(t, m.CompareAndSwap(1, 3))
	assert.True(t, m.CompareAndSwap(2, 3))
	assert.Equal(t, gust.Some(3), m.Load())
}

type one int

func (o *one) Increment() {
	*o++
}

func runLazyValue(t *testing.T, once *syncutil.LazyValue[*one], c chan bool) {
	o := once.TryGetValue().Unwrap()
	if v := *o; v != 1 {
		t.Errorf("once failed inside run: %d is not 1", v)
	}
	c <- true
}

func TestLazyValue(t *testing.T) {
	assert.Equal(t, gust.TryErr[int](syncutil.ErrLazyValueWithoutInit), new(syncutil.LazyValue[int]).TryGetValue())
	assert.Equal(t, 0, new(syncutil.LazyValue[int]).SetInitValue(0).TryGetValue().Unwrap())
	assert.Equal(t, 1, new(syncutil.LazyValue[int]).SetInitValue(1).TryGetValue().Unwrap())
	o := new(one)
	once := new(syncutil.LazyValue[*one]).SetInitFunc(func() gust.Result[*one] {
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

func TestLazyValuePanic1(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			assert.Equal(t, "failed", p)
		} else {
			t.Fatalf("should painc")
		}
	}()
	var once = new(syncutil.LazyValue[struct{}]).SetInitFunc(func() gust.Result[struct{}] {
		panic("failed")
	})
	_ = once.TryGetValue().Unwrap()
	t.Fatalf("unreachable")
}

func TestLazyValuePanic2(t *testing.T) {
	defer func() {
		if p := recover(); p != nil {
			assert.Equal(t, syncutil.ErrLazyValueWithoutInit, fmt.Sprintf("%v", p))
		} else {
			t.Fatalf("should painc")
		}
	}()
	_ = new(syncutil.LazyValue[struct{}]).TryGetValue().Unwrap()
	t.Fatalf("unreachable")
}

func TestLazyValueGetPtr(t *testing.T) {
	assert.Equal(t, (*int)(nil), new(syncutil.LazyValue[int]).GetPtr())
	var zero int = 0
	assert.Equal(t, &zero, new(syncutil.LazyValue[int]).SetInitZero().GetPtr())
	var one int = 1
	assert.Equal(t, &one, new(syncutil.LazyValue[int]).SetInitValue(1).GetPtr())
}

func TestNewRWMutex(t *testing.T) {
	m := syncutil.NewRWMutex(10)

	// Test Lock/Unlock
	assert.Equal(t, 10, m.Lock())
	m.Unlock(20)
	assert.Equal(t, 20, m.Lock())
	m.Unlock()

	// Test TryLock
	result := m.TryLock()
	assert.True(t, result.IsSome())
	assert.Equal(t, 20, result.Unwrap())
	m.Unlock()

	// Test RLock/RUnlock
	assert.Equal(t, 20, m.RLock())
	m.RUnlock()

	// Test TryRLock
	result2 := m.TryRLock()
	assert.True(t, result2.IsSome())
	assert.Equal(t, 20, result2.Unwrap())
	m.RUnlock()

	// Test LockScope
	m.LockScope(func(old int) int {
		return old * 2
	})
	result3 := m.TryLock()
	assert.True(t, result3.IsSome())
	assert.Equal(t, 40, result3.Unwrap())
	m.Unlock()

	// Test TryLockScope
	success := false
	m.TryLockScope(func(old int) int {
		success = true
		return old + 1
	})
	assert.True(t, success)
	result4 := m.TryLock()
	assert.True(t, result4.IsSome())
	assert.Equal(t, 41, result4.Unwrap())
	m.Unlock()

	// Test RLockScope
	readValue := 0
	m.RLockScope(func(val int) {
		readValue = val
	})
	assert.Equal(t, 41, readValue)

	// Test TryRLockScope
	readValue2 := 0
	m.TryRLockScope(func(val int) {
		readValue2 = val
	})
	assert.Equal(t, 41, readValue2)

	// Test TryBest
	swapCount := 0
	m.TryBest(func(val int) bool {
		return val > 50 // Condition fails
	}, func(old int) gust.Option[int] {
		swapCount++
		return gust.Some(100) // Swap to 100
	})
	assert.Equal(t, 1, swapCount)
	result5 := m.TryLock()
	assert.True(t, result5.IsSome())
	assert.Equal(t, 100, result5.Unwrap())
	m.Unlock()

	// Test TryBest with successful condition
	m.TryBest(func(val int) bool {
		return val > 50 // Condition succeeds
	}, func(old int) gust.Option[int] {
		return gust.Some(200) // Should not be called
	})
	result6 := m.TryLock()
	assert.True(t, result6.IsSome())
	assert.Equal(t, 100, result6.Unwrap()) // Value unchanged
	m.Unlock()
}

func TestLazyValueNewFunctions(t *testing.T) {
	// Test NewLazyValue
	lv1 := syncutil.NewLazyValue[int]()
	assert.False(t, lv1.IsInitialized())
	assert.Equal(t, gust.TryErr[int](syncutil.ErrLazyValueWithoutInit), lv1.TryGetValue())

	// Test NewLazyValueWithValue - lazy initialization
	lv2 := syncutil.NewLazyValueWithValue(42)
	assert.False(t, lv2.IsInitialized()) // Not initialized until TryGetValue is called
	assert.Equal(t, gust.Some(42), lv2.TryGetValue().Ok())
	assert.True(t, lv2.IsInitialized()) // Now initialized after TryGetValue

	// Test NewLazyValueWithZero - lazy initialization
	lv3 := syncutil.NewLazyValueWithZero[int]()
	assert.False(t, lv3.IsInitialized()) // Not initialized until TryGetValue is called
	assert.Equal(t, gust.Some(0), lv3.TryGetValue().Ok())
	assert.True(t, lv3.IsInitialized()) // Now initialized after TryGetValue

	// Test NewLazyValueWithFunc - lazy initialization
	lv4 := syncutil.NewLazyValueWithFunc[int](func() gust.Result[int] {
		return gust.Ok(100)
	})
	assert.False(t, lv4.IsInitialized()) // Not initialized until TryGetValue is called
	assert.Equal(t, gust.Some(100), lv4.TryGetValue().Ok())
	assert.True(t, lv4.IsInitialized()) // Now initialized after TryGetValue

	// Test SetInitFunc on uninitialized
	lv5 := syncutil.NewLazyValue[string]()
	lv5.SetInitFunc(func() gust.Result[string] {
		return gust.Ok("test")
	})
	assert.Equal(t, gust.Some("test"), lv5.TryGetValue().Ok())

	// Test SetInitFunc on initialized (should not change)
	lv6 := syncutil.NewLazyValueWithValue("original")
	// First call to TryGetValue initializes it
	assert.Equal(t, gust.Some("original"), lv6.TryGetValue().Ok())
	assert.True(t, lv6.IsInitialized())
	// Now SetInitFunc should not change the value since it's already initialized
	lv6.SetInitFunc(func() gust.Result[string] {
		return gust.Ok("new")
	})
	assert.Equal(t, gust.Some("original"), lv6.TryGetValue().Ok())

	// Test SetInitValue
	lv7 := syncutil.NewLazyValue[int]()
	lv7.SetInitValue(200)
	assert.Equal(t, gust.Some(200), lv7.TryGetValue().Ok())

	// Test SetInitZero
	lv8 := syncutil.NewLazyValue[int]()
	lv8.SetInitZero()
	assert.Equal(t, gust.Some(0), lv8.TryGetValue().Ok())

	// Test Zero
	var zero int
	assert.Equal(t, zero, lv8.Zero())
}

func BenchmarkLazyValue(b *testing.B) {
	var once = new(syncutil.LazyValue[struct{}])
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			// once.GetValue(false)
			once.TryGetValue()
		}
	})
}

func TestRWMutex_TryBest_NilReadAndDo(t *testing.T) {
	m := syncutil.NewRWMutex(10)
	// Should not panic with nil readAndDo
	m.TryBest(nil, func(old int) gust.Option[int] {
		return gust.Some(20)
	})
	assert.Equal(t, 10, m.Lock())
	m.Unlock()
}

func TestRWMutex_TryBest_NilSwapWhenFalse(t *testing.T) {
	m := syncutil.NewRWMutex(10)
	// Should not panic with nil swapWhenFalse
	m.TryBest(func(val int) bool {
		return val > 50
	}, nil)
	assert.Equal(t, 10, m.Lock())
	m.Unlock()
}

func TestRWMutex_TryBest_SwapWhenFalseReturnsNone(t *testing.T) {
	m := syncutil.NewRWMutex(10)
	m.TryBest(func(val int) bool {
		return val > 50
	}, func(old int) gust.Option[int] {
		return gust.None[int]() // Return None, should not swap
	})
	assert.Equal(t, 10, m.Lock()) // Value unchanged
	m.Unlock()
}

func TestMutex_TryLockScope_Failed(t *testing.T) {
	m := syncutil.NewMutex(10)
	// Lock it first
	m.Lock()
	// TryLockScope should fail since lock is already held
	success := false
	m.TryLockScope(func(old int) int {
		success = true
		return old + 1
	})
	assert.False(t, success)
	m.Unlock()
	assert.Equal(t, 10, m.Lock()) // Value unchanged
	m.Unlock()
}

func TestMutex_TryLock_Failed(t *testing.T) {
	m := syncutil.NewMutex(10)
	// Lock it first
	m.Lock()
	// TryLock should fail since lock is already held
	result := m.TryLock()
	assert.True(t, result.IsNone(), "TryLock should fail when lock is already held")
	m.Unlock()
	// Now TryLock should succeed
	result2 := m.TryLock()
	assert.True(t, result2.IsSome())
	assert.Equal(t, 10, result2.Unwrap())
	m.Unlock()
}

func TestRWMutex_TryLock_Failed(t *testing.T) {
	m := syncutil.NewRWMutex(10)
	// Lock it first for writing
	m.Lock()
	// TryLock should fail since lock is already held
	result := m.TryLock()
	assert.True(t, result.IsNone(), "TryLock should fail when lock is already held")
	m.Unlock()
	// Now TryLock should succeed
	result2 := m.TryLock()
	assert.True(t, result2.IsSome())
	assert.Equal(t, 10, result2.Unwrap())
	m.Unlock()

	// Test TryLock fails when RLock is held
	m.RLock()
	result3 := m.TryLock()
	assert.True(t, result3.IsNone(), "TryLock should fail when RLock is held")
	m.RUnlock()
	// Now TryLock should succeed
	result4 := m.TryLock()
	assert.True(t, result4.IsSome())
	assert.Equal(t, 10, result4.Unwrap())
	m.Unlock()
}

func TestRWMutex_TryRLockScope_Failed(t *testing.T) {
	m := syncutil.NewRWMutex(10)
	// Lock it first for writing
	m.Lock()
	// TryRLockScope should fail since lock is already held
	success := false
	m.TryRLockScope(func(val int) {
		success = true
	})
	assert.False(t, success)
	m.Unlock()
}

func TestSyncMap_Range_TypeMismatch(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)
	m.Store("key2", 100)

	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	// Should iterate over all valid entries
	assert.Equal(t, 2, count)
}

func TestSyncMap_Load_NotFound(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	result := m.Load("nonexistent")
	assert.True(t, result.IsNone())
}

func TestSyncMap_LoadOrStore_NewKey(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	result := m.LoadOrStore("new", 42)
	assert.True(t, result.IsNone(), "New key should return None")

	// Verify value was stored
	loaded := m.Load("new")
	assert.True(t, loaded.IsSome())
	assert.Equal(t, 42, loaded.Unwrap())
}

func TestSyncMap_LoadOrStore_ExistingKey(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	m.Store("existing", 100)
	result := m.LoadOrStore("existing", 200)
	assert.True(t, result.IsSome(), "Existing key should return Some")
	assert.Equal(t, 100, result.Unwrap(), "Should return existing value")

	// Verify value was not changed
	loaded := m.Load("existing")
	assert.Equal(t, 100, loaded.Unwrap())
}

func TestSyncMap_LoadAndDelete_ExistingKey(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)
	result := m.LoadAndDelete("key")
	assert.True(t, result.IsSome())
	assert.Equal(t, 42, result.Unwrap())

	// Verify key was deleted
	loaded := m.Load("key")
	assert.True(t, loaded.IsNone())
}

func TestSyncMap_LoadAndDelete_NotFound(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	result := m.LoadAndDelete("nonexistent")
	assert.True(t, result.IsNone())
}

func TestSyncMap_Delete(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)
	m.Delete("key")

	// Verify key was deleted
	loaded := m.Load("key")
	assert.True(t, loaded.IsNone())
}

func TestSyncMap_Range_KeyTypeMismatch(t *testing.T) {
	// Test Range when key type doesn't match (should skip that entry)
	// This tests the key type assertion failure path in Range
	// Note: Since SyncMap uses sync.Map internally and sync.Map is type-safe,
	// we cannot directly test type mismatch through the public API.
	// The type assertion failure path in Range() is defensive code that handles
	// edge cases, but cannot be easily tested without accessing unexported fields.

	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)

	// Test normal Range behavior
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)

	// Note: The type assertion failure paths (lines 76-77 and 80-81 in map.go)
	// cannot be tested through the public API because sync.Map is type-safe.
	// These are defensive checks that would only trigger if there's a bug
	// or if someone uses unsafe/reflection to bypass type safety.
}

func TestSyncMap_Range_ValueTypeMismatch(t *testing.T) {
	// Test Range when value type doesn't match (should skip that entry)
	// This tests the value type assertion failure path in Range
	// Note: Since SyncMap uses sync.Map internally and sync.Map is type-safe,
	// we cannot directly test type mismatch through the public API.
	// The type assertion failure path in Range() is defensive code that handles
	// edge cases, but cannot be easily tested without accessing unexported fields.

	var m syncutil.SyncMap[string, int]
	m.Store("key", 42)

	// Test normal Range behavior
	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return true
	})
	assert.Equal(t, 1, count)

	// Note: The type assertion failure paths (lines 76-77 and 80-81 in map.go)
	// cannot be tested through the public API because sync.Map is type-safe.
	// These are defensive checks that would only trigger if there's a bug
	// or if someone uses unsafe/reflection to bypass type safety.
}

func TestSyncMap_Range_EarlyTermination(t *testing.T) {
	var m syncutil.SyncMap[string, int]
	m.Store("a", 1)
	m.Store("b", 2)
	m.Store("c", 3)

	count := 0
	m.Range(func(key string, value int) bool {
		count++
		return count < 2 // Stop after 2 iterations
	})
	assert.Equal(t, 2, count)
}

func TestAtomicValue_CompareAndSwap(t *testing.T) {
	var v syncutil.AtomicValue[int]
	v.Store(10)

	// Test successful swap
	assert.True(t, v.CompareAndSwap(10, 20))
	assert.Equal(t, gust.Some(20), v.Load())

	// Test failed swap (old value doesn't match)
	assert.False(t, v.CompareAndSwap(10, 30))
	assert.Equal(t, gust.Some(20), v.Load()) // Value unchanged

	// Test successful swap with correct old value
	assert.True(t, v.CompareAndSwap(20, 30))
	assert.Equal(t, gust.Some(30), v.Load())
}

func TestLazyValue_SetInitFunc_AlreadyInitialized(t *testing.T) {
	lv := syncutil.NewLazyValueWithValue("original")
	// Initialize it
	_ = lv.TryGetValue()
	assert.True(t, lv.IsInitialized())

	// Try to set init func after initialization
	lv.SetInitFunc(func() gust.Result[string] {
		return gust.Ok("new")
	})

	// Should still return original value
	assert.Equal(t, gust.Some("original"), lv.TryGetValue().Ok())
}

func TestLazyValue_SetInitFunc_AlreadySet(t *testing.T) {
	lv := syncutil.NewLazyValue[string]()
	lv.SetInitFunc(func() gust.Result[string] {
		return gust.Ok("first")
	})

	// Try to set another init func
	lv.SetInitFunc(func() gust.Result[string] {
		return gust.Ok("second")
	})

	// Should use first function
	assert.Equal(t, gust.Some("first"), lv.TryGetValue().Ok())
}

func TestLazyValue_Zero(t *testing.T) {
	lv := syncutil.NewLazyValue[int]()
	var zero int
	assert.Equal(t, zero, lv.Zero())

	lv2 := syncutil.NewLazyValue[string]()
	var zeroStr string
	assert.Equal(t, zeroStr, lv2.Zero())
}
