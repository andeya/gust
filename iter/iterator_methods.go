package iter

import (
	stditer "iter"

	"github.com/andeya/gust"
)

// MustToDoubleEnded converts to a DoubleEndedIterator[T] if the underlying
// iterator supports double-ended iteration. Otherwise, it panics.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var deIter = iter.MustToDoubleEnded()
//	assert.Equal(t, gust.Some(3), deIter.NextBack())
//	// Can use Iterator methods:
//	var doubled = deIter.Map(func(x int) any { return x * 2 })
func (it Iterator[T]) MustToDoubleEnded() DoubleEndedIterator[T] {
	if deCore, ok := it.iterable.(DoubleEndedIterable[T]); ok {
		return DoubleEndedIterator[T]{
			Iterator: Iterator[T]{iterable: deCore}, // Embed Iterator with the same core
			iterable: deCore,
		}
	}
	panic("iterator does not support double-ended iteration")
}

// TryToDoubleEnded converts to a DoubleEndedIterator[T] if the underlying
// iterator supports double-ended iteration. Otherwise, it returns None.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var deIter = iter.TryToDoubleEnded()
//	assert.Equal(t, gust.Some(3), deIter.NextBack())
//	// Can use Iterator methods:
//	var doubled = deIter.Map(func(x int) any { return x * 2 })
func (it Iterator[T]) TryToDoubleEnded() gust.Option[DoubleEndedIterator[T]] {
	if deCore, ok := it.iterable.(DoubleEndedIterable[T]); ok {
		return gust.Some(DoubleEndedIterator[T]{
			Iterator: Iterator[T]{iterable: deCore}, // Embed Iterator with the same core
			iterable: deCore,
		})
	}
	return gust.None[DoubleEndedIterator[T]]()
}

// Next advances the iterator and returns the next value.
// This implements gust.Iterable[T] interface.
//
//go:inline
func (it Iterator[T]) Next() gust.Option[T] {
	return it.iterable.Next()
}

// SizeHint returns the bounds on the remaining length of the iterator.
// This implements gust.IterableSizeHint interface.
//
//go:inline
func (it Iterator[T]) SizeHint() (uint, gust.Option[uint]) {
	return it.iterable.SizeHint()
}

// Seq converts the Iterator[T] to Go's standard iter.Seq[T].
// This allows using gust iterators with Go's built-in iteration support (for loops).
//
// # Examples
//
//	iter := FromSlice([]int{1, 2, 3})
//	for v := range iter.Seq() {
//		fmt.Println(v) // prints 1, 2, 3
//	}
//
//	// Works with Go's standard library functions
//	iter := FromSlice([]int{1, 2, 3})
//	all := iter.Seq().All(func(v int) bool { return v > 0 })
func (it Iterator[T]) Seq() stditer.Seq[T] {
	return func(yield func(T) bool) {
		for {
			opt := it.Next()
			if opt.IsNone() {
				return
			}
			if !yield(opt.Unwrap()) {
				return
			}
		}
	}
}

// Pull converts the Iterator[T] to a pull-style iterator using Go's standard iter.Pull.
// This returns two functions: next (to pull values) and stop (to stop iteration).
// The caller should defer stop() to ensure proper cleanup.
//
// # Examples
//
//	iter := FromSlice([]int{1, 2, 3, 4, 5})
//	next, stop := iter.Pull()
//	defer stop()
//
//	// Pull values manually
//	for {
//		v, ok := next()
//		if !ok {
//			break
//		}
//		fmt.Println(v)
//		if v == 3 {
//			break // Early termination
//		}
//	}
func (it Iterator[T]) Pull() (next func() (T, bool), stop func()) {
	return stditer.Pull(it.Seq())
}

// Adapter methods - these return new iterators and can be chained
// Note: Methods that change the type (like Map) must use function-style API
// due to Go's limitation on generic methods in structs.

// Filter creates an iterator which uses a closure to determine if an element should be yielded.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2})
//	var filtered = iter.Filter(func(x int) bool { return x > 0 })
//	assert.Equal(t, gust.Some(1), filtered.Next())
//
//go:inline
func (it Iterator[T]) Filter(predicate func(T) bool) Iterator[T] {
	return filterImpl(it, predicate)
}

// Chain takes two iterators and creates a new iterator over both in sequence.
//
// # Examples
//
//	var iter1 = FromSlice([]int{1, 2, 3})
//	var iter2 = FromSlice([]int{4, 5, 6})
//	var chained = iter1.Chain(iter2)
//	assert.Equal(t, gust.Some(1), chained.Next())
//
//go:inline
func (it Iterator[T]) Chain(other Iterator[T]) Iterator[T] {
	return chainImpl(it, other)
}

// Skip creates an iterator that skips the first n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var skipped = iter.Skip(2)
//	assert.Equal(t, gust.Some(3), skipped.Next())
//
//go:inline
func (it Iterator[T]) Skip(n uint) Iterator[T] {
	return skipImpl(it, n)
}

// Take creates an iterator that yields the first n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var taken = iter.Take(2)
//	assert.Equal(t, gust.Some(1), taken.Next())
//	assert.Equal(t, gust.Some(2), taken.Next())
//	assert.Equal(t, gust.None[int](), taken.Next())
//
//go:inline
func (it Iterator[T]) Take(n uint) Iterator[T] {
	return takeImpl(it, n)
}

// StepBy creates an iterator that steps by n elements at a time.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2, 3, 4, 5})
//	var stepped = iter.StepBy(2)
//	assert.Equal(t, gust.Some(0), stepped.Next())
//	assert.Equal(t, gust.Some(2), stepped.Next())
//	assert.Equal(t, gust.Some(4), stepped.Next())
//
//go:inline
func (it Iterator[T]) StepBy(step uint) Iterator[T] {
	return stepByImpl(it, step)
}

// SkipWhile creates an iterator that skips elements while a predicate is true.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var skipped = iter.SkipWhile(func(x int) bool { return x < 3 })
//	assert.Equal(t, gust.Some(3), skipped.Next())
//
//go:inline
func (it Iterator[T]) SkipWhile(predicate func(T) bool) Iterator[T] {
	return skipWhileImpl(it, predicate)
}

// TakeWhile creates an iterator that yields elements while a predicate is true.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var taken = iter.TakeWhile(func(x int) bool { return x < 3 })
//	assert.Equal(t, gust.Some(1), taken.Next())
//	assert.Equal(t, gust.Some(2), taken.Next())
//	assert.Equal(t, gust.None[int](), taken.Next())
//
//go:inline
func (it Iterator[T]) TakeWhile(predicate func(T) bool) Iterator[T] {
	return takeWhileImpl(it, predicate)
}

// Fuse creates an iterator that ends after the first None.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var fused = iter.Fuse()
//	// After None, it will always return None
//
//go:inline
func (it Iterator[T]) Fuse() Iterator[T] {
	return fuseImpl(it)
}

// Inspect creates an iterator that calls a closure on each element for side effects.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var inspected = iter.Inspect(func(x int) { fmt.Println(x) })
//	// Prints 1, 2, 3 when iterated
//
//go:inline
func (it Iterator[T]) Inspect(f func(T)) Iterator[T] {
	return inspectImpl(it, f)
}

// Consumer methods - these consume the iterator

// Collect collects all items into a slice.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var result = iter.Collect()
//	assert.Equal(t, []int{1, 2, 3}, result)
//
//go:inline
func (it Iterator[T]) Collect() []T {
	return collectImpl(it.iterable)
}

// Count consumes the iterator, counting the number of iterations.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	assert.Equal(t, uint(5), iter.Count())
//
//go:inline
func (it Iterator[T]) Count() uint {
	return countImpl(it.iterable)
}

// Last returns the last element of the iterator.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, gust.Some(3), iter.Last())
//
//go:inline
func (it Iterator[T]) Last() gust.Option[T] {
	return lastImpl(it.iterable)
}

// Reduce reduces the iterator to a single value.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iter.Reduce(func(acc int, x int) int { return acc + x })
//	assert.Equal(t, gust.Some(6), sum)
//
//go:inline
func (it Iterator[T]) Reduce(f func(T, T) T) gust.Option[T] {
	return reduceImpl(it.iterable, f)
}

// ForEach calls a closure on each element.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	iter.ForEach(func(x int) { fmt.Println(x) })
//
//go:inline
func (it Iterator[T]) ForEach(f func(T)) {
	forEachImpl(it.iterable, f)
}

// All tests if all elements satisfy a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{2, 4, 6})
//	assert.True(t, iter.All(func(x int) bool { return x%2 == 0 }))
//
//go:inline
func (it Iterator[T]) All(predicate func(T) bool) bool {
	return allImpl(it.iterable, predicate)
}

// Any tests if any element satisfies a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.True(t, iter.Any(func(x int) bool { return x > 2 }))
//
//go:inline
func (it Iterator[T]) Any(predicate func(T) bool) bool {
	return anyImpl(it.iterable, predicate)
}

// Find searches for an element that satisfies a predicate.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, gust.Some(2), iter.Find(func(x int) bool { return x > 1 }))
//
//go:inline
func (it Iterator[T]) Find(predicate func(T) bool) gust.Option[T] {
	return findImpl(it.iterable, predicate)
}

// XFindMap searches for an element and maps it (any version).
//
// # Examples
//
//	var iter = FromSlice([]string{"lol", "NaN", "2", "5"})
//	var firstNumber = iter.XFindMap(func(s string) gust.Option[any] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Some(any(v))
//		}
//		return gust.None[any]()
//	})
//	assert.True(t, firstNumber.IsSome())
//	assert.Equal(t, 2, firstNumber.Unwrap().(int))
//
//go:inline
func (it Iterator[T]) XFindMap(f func(T) gust.Option[any]) gust.Option[any] {
	return FindMap(it, f)
}

// FindMap searches for an element and maps it.
//
// # Examples
//
//	var iter = FromSlice([]string{"lol", "NaN", "2", "5"})
//	var firstNumber = iter.FindMap(func(s string) gust.Option[int] {
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Some(v)
//		}
//		return gust.None[int]()
//	})
//	assert.Equal(t, gust.Some(2), firstNumber)
//
//go:inline
func (it Iterator[T]) FindMap(f func(T) gust.Option[T]) gust.Option[T] {
	return FindMap(it, f)
}

// XMap creates an iterator which calls a closure on each element (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var doubled = iter.XMap(func(x int) any { return x * 2 })
//	assert.Equal(t, gust.Some(2), doubled.Next())
//	// Can chain: doubled.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XMap(f func(T) any) Iterator[any] {
	return Map(it, f)
}

// Map creates an iterator which calls a closure on each element.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var doubled = iter.Map(func(x int) int { return x * 2 })
//	assert.Equal(t, gust.Some(2), doubled.Next())
//	// Can chain: doubled.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) Map(f func(T) T) Iterator[T] {
	return Map(it, f)
}

// XFlatMap creates an iterator that maps and flattens nested iterators (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var flatMapped = iter.XFlatMap(func(x int) Iterator[any] {
//		return FromSlice([]any{x, x * 2})
//	})
//	// Can chain: flatMapped.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XFlatMap(f func(T) Iterator[any]) Iterator[any] {
	return FlatMap(it, f)
}

// FlatMap creates an iterator that maps and flattens nested iterators.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var flatMapped = iter.FlatMap(func(x int) Iterator[int] {
//		return FromSlice([]int{x, x * 2})
//	})
//	// Can chain: flatMapped.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) FlatMap(f func(T) Iterator[T]) Iterator[T] {
	return FlatMap(it, f)
}

// XFold folds every element into an accumulator.
// This wrapper method allows XFold to be called as a method.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iter.XFold(0, func(acc any, x int) any { return acc.(int) + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (it Iterator[T]) XFold(init any, f func(any, T) any) any {
	return Fold(it, init, f)
}

// Fold folds every element into an accumulator.
// This wrapper method allows Fold to be called as a method.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iter.Fold(0, func(acc int, x int) int { return acc + x })
//	assert.Equal(t, 6, sum)
//
//go:inline
func (it Iterator[T]) Fold(init T, f func(T, T) T) T {
	return Fold(it, init, f)
}

// XTryFold applies a function as long as it returns successfully, producing a single, final value (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iter.XTryFold(0, func(acc any, x int) gust.Result[any] {
//		return gust.Ok(any(acc.(int) + x))
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap().(int))
//
//go:inline
func (it Iterator[T]) XTryFold(init any, f func(any, T) gust.Result[any]) gust.Result[any] {
	return TryFold(it, init, f)
}

// TryFold applies a function as long as it returns successfully, producing a single, final value.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var sum = iter.TryFold(0, func(acc int, x int) gust.Result[int] {
//		return gust.Ok(acc + x)
//	})
//	assert.True(t, sum.IsOk())
//	assert.Equal(t, 6, sum.Unwrap())
//
//go:inline
func (it Iterator[T]) TryFold(init T, f func(T, T) gust.Result[T]) gust.Result[T] {
	return TryFold(it, init, f)
}

// Partition partitions the iterator into two slices.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	evens, odds := iter.Partition(func(x int) bool { return x%2 == 0 })
//	assert.Equal(t, []int{2, 4}, evens)
//	assert.Equal(t, []int{1, 3, 5}, odds)
//
//go:inline
func (it Iterator[T]) Partition(f func(T) bool) (truePart []T, falsePart []T) {
	return partitionImpl(it.iterable, f)
}

// Position searches for an element in an iterator, returning its index.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, gust.Some(uint(1)), iter.Position(func(x int) bool { return x == 2 }))
//
//go:inline
//go:inline
func (it Iterator[T]) Position(predicate func(T) bool) gust.Option[uint] {
	return positionImpl(it.iterable, predicate)
}

// AdvanceBy advances the iterator by n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4})
//	assert.Equal(t, gust.NonErrable[uint](), iter.AdvanceBy(2))
//
//go:inline
func (it Iterator[T]) AdvanceBy(n uint) gust.Errable[uint] {
	return advanceByImpl(it.iterable, n)
}

// Nth returns the nth element of the iterator.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	assert.Equal(t, gust.Some(2), iter.Nth(1))
//
//go:inline
func (it Iterator[T]) Nth(n uint) gust.Option[T] {
	return nthImpl(it.iterable, n)
}

// NextChunk advances the iterator and returns an array containing the next N values.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	chunk := iter.NextChunk(2)
//	assert.True(t, chunk.IsOk())
//
//go:inline
func (it Iterator[T]) NextChunk(n uint) gust.EnumResult[[]T, []T] {
	return nextChunkImpl(it.iterable, n)
}

// Intersperse creates an iterator that places a separator between adjacent items.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2})
//	var interspersed = iter.Intersperse(100)
//	assert.Equal(t, gust.Some(0), interspersed.Next())
//	assert.Equal(t, gust.Some(100), interspersed.Next())
//
//go:inline
func (it Iterator[T]) Intersperse(separator T) Iterator[T] {
	return intersperseImpl(it.iterable, separator)
}

// IntersperseWith creates an iterator that places an item generated by separator between adjacent items.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2})
//	var interspersed = iter.IntersperseWith(func() int { return 99 })
//	assert.Equal(t, gust.Some(0), interspersed.Next())
//	assert.Equal(t, gust.Some(99), interspersed.Next())
//
//go:inline
func (it Iterator[T]) IntersperseWith(separator func() T) Iterator[T] {
	return intersperseWithImpl(it.iterable, separator)
}

// Cycle repeats an iterator endlessly.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var cycled = iter.Cycle()
//	assert.Equal(t, gust.Some(1), cycled.Next())
//	assert.Equal(t, gust.Some(2), cycled.Next())
//	assert.Equal(t, gust.Some(3), cycled.Next())
//	assert.Equal(t, gust.Some(1), cycled.Next()) // starts over
//
//go:inline
func (it Iterator[T]) Cycle() Iterator[T] {
	return cycleImpl(it.iterable)
}

// XFilterMap creates an iterator that both filters and maps (any version).
//
// # Examples
//
//	var iter = FromSlice([]string{"1", "two", "NaN", "four", "5"})
//	var filtered = iter.XFilterMap(func(s string) gust.Option[any] {
//		if s == "1" {
//			return gust.Some(any(1))
//		}
//		if s == "5" {
//			return gust.Some(any(5))
//		}
//		return gust.None[any]()
//	})
//	// Can chain: filtered.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XFilterMap(f func(T) gust.Option[any]) Iterator[any] {
	return FilterMap(it, f)
}

// FilterMap creates an iterator that both filters and maps.
//
// # Examples
//
//	var iter = FromSlice([]string{"1", "two", "NaN", "four", "5"})
//	var filtered = iter.FilterMap(func(s string) gust.Option[string] {
//		if s == "1" || s == "5" {
//			return gust.Some(s)
//		}
//		return gust.None[string]()
//	})
//	// Can chain: filtered.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) FilterMap(f func(T) gust.Option[T]) Iterator[T] {
	return FilterMap(it, f)
}

// XMapWhile creates an iterator that both yields elements based on a predicate and maps (any version).
//
//go:inline
func (it Iterator[T]) XMapWhile(predicate func(T) gust.Option[any]) Iterator[any] {
	return MapWhile(it, predicate)
}

// MapWhile creates an iterator that both yields elements based on a predicate and maps.
//
//go:inline
func (it Iterator[T]) MapWhile(predicate func(T) gust.Option[T]) Iterator[T] {
	return MapWhile(it, predicate)
}

// XScan creates an iterator that scans over the iterator with a state (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var scanned = iter.XScan(0, func(state *any, x int) gust.Option[any] {
//		s := (*state).(int) + x
//		*state = s
//		return gust.Some(any(s))
//	})
//	// Can chain: scanned.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XScan(initialState any, f func(*any, T) gust.Option[any]) Iterator[any] {
	return Scan(it, initialState, f)
}

// Scan creates an iterator that scans over the iterator with a state.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var scanned = iter.Scan(0, func(state *int, x int) gust.Option[int] {
//		*state = *state + x
//		return gust.Some(*state)
//	})
//	// Can chain: scanned.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) Scan(initialState T, f func(*T, T) gust.Option[T]) Iterator[T] {
	return Scan(it, initialState, f)
}

// XMapWindows creates an iterator that applies a function to overlapping windows (any version).
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var windows = iter.XMapWindows(3, func(window []int) any {
//		return window[0] + window[1] + window[2]
//	})
//	// Can chain: windows.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) XMapWindows(windowSize uint, f func([]T) any) Iterator[any] {
	return MapWindows(it, windowSize, f)
}

// MapWindows creates an iterator that applies a function to overlapping windows.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3, 4, 5})
//	var windows = iter.MapWindows(3, func(window []int) int {
//		return window[0] + window[1] + window[2]
//	})
//	// Can chain: windows.Filter(...).Collect()
//
//go:inline
func (it Iterator[T]) MapWindows(windowSize uint, f func([]T) T) Iterator[T] {
	return MapWindows(it, windowSize, f)
}

// MaxBy returns the element that gives the maximum value with respect to the
// specified comparison function.
//
// # Examples
//
//	var iter = FromSlice([]int{-3, 0, 1, 5, -10})
//	var max = iter.MaxBy(func(x, y int) int {
//		if x < y {
//			return -1
//		}
//		if x > y {
//			return 1
//		}
//		return 0
//	})
//	assert.Equal(t, gust.Some(5), max)
//
//go:inline
func (it Iterator[T]) MaxBy(compare func(T, T) int) gust.Option[T] {
	return maxByImpl(it, compare)
}

// MinBy returns the element that gives the minimum value with respect to the
// specified comparison function.
//
// # Examples
//
//	var iter = FromSlice([]int{-3, 0, 1, 5, -10})
//	var min = iter.MinBy(func(x, y int) int {
//		if x < y {
//			return -1
//		}
//		if x > y {
//			return 1
//		}
//		return 0
//	})
//	assert.Equal(t, gust.Some(-10), min)
//
//go:inline
func (it Iterator[T]) MinBy(compare func(T, T) int) gust.Option[T] {
	return minByImpl(it, compare)
}

// XTryForEach applies a fallible function to each item in the iterator,
// stopping at the first error and returning that error.
//
// # Examples
//
//	var data = []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
//	var res = iter.TryForEach(func(x string) gust.Result[any] {
//		fmt.Println(x)
//		return gust.Ok[any](nil)
//	})
//	assert.True(t, res.IsOk())
//
//go:inline
func (it Iterator[T]) XTryForEach(f func(T) gust.Result[any]) gust.Result[any] {
	return TryForEach(it, f)
}

// TryForEach applies a fallible function to each item in the iterator,
// stopping at the first error and returning that error.
//
// # Examples
//
//	var data = []string{"no_tea.txt", "stale_bread.json", "torrential_rain.png"}
//	var res = iter.TryForEach(func(x string) gust.Result[string] {
//		fmt.Println(x)
//		return gust.Ok[string](x+"_processed")
//	})
//	assert.True(t, res.IsOk())
//	assert.Equal(t, "no_tea.txt_processed", res.Unwrap())
//
//go:inline
func (it Iterator[T]) TryForEach(f func(T) gust.Result[T]) gust.Result[T] {
	return TryForEach(it, f)
}

// TryReduce reduces the elements to a single one by repeatedly applying a reducing operation.
//
// # Examples
//
//	var numbers = []int{10, 20, 5, 23, 0}
//	var sum = iter.TryReduce(func(x, y int) gust.Result[int] {
//		if x+y > 100 {
//			return gust.Err[int](errors.New("overflow"))
//		}
//		return gust.Ok(x + y)
//	})
//	assert.True(t, sum.IsOk())
//
//go:inline
func (it Iterator[T]) TryReduce(f func(T, T) gust.Result[T]) gust.Result[gust.Option[T]] {
	return tryReduceImpl(it, f)
}

// TryFind applies function to the elements of iterator and returns
// the first true result or the first error.
//
// # Examples
//
//	var a = []string{"1", "2", "lol", "NaN", "5"}
//	var result = iter.TryFind(func(s string) gust.Result[bool] {
//		if s == "lol" {
//			return gust.Err[bool](errors.New("invalid"))
//		}
//		if v, err := strconv.Atoi(s); err == nil {
//			return gust.Ok(v == 2)
//		}
//		return gust.Ok(false)
//	})
//	assert.True(t, result.IsOk())
//
//go:inline
func (it Iterator[T]) TryFind(f func(T) gust.Result[bool]) gust.Result[gust.Option[T]] {
	return tryFindImpl(it, f)
}

// Peekable creates a peekable iterator.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var peekable = iter.Peekable()
//	assert.Equal(t, gust.Some(1), peekable.Peek())
//	assert.Equal(t, gust.Some(1), peekable.Next())
//
//	// Can use all Iterator methods:
//	var filtered = peekable.Filter(func(x int) bool { return x > 1 })
//	assert.Equal(t, gust.Some(2), filtered.Next())
//
//go:inline
func (it Iterator[T]) Peekable() PeekableIterator[T] {
	return peekableImpl(it)
}
