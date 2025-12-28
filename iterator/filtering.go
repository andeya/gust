package iterator

import (
	"github.com/andeya/gust/option"
	"github.com/andeya/gust/result"
	"github.com/andeya/gust/void"
)

// advanceByImplFiltering advances an iterator by n elements.
// This is a helper function used internally by Skip and StepBy in this module.
//
//go:inline
func advanceByImplFiltering[T any](iter Iterable[T], n uint) result.VoidResult {
	it := Iterator[T]{iterable: iter}
	for i := uint(0); i < n; i++ {
		if it.Next().IsNone() {
			return result.TryErr[void.Void](n - i)
		}
	}
	return result.Ok[void.Void](nil)
}

//go:inline
func skipImpl[T any](iter Iterator[T], n uint) Iterator[T] {
	return Iterator[T]{iterable: &skipIterable[T]{iter: iter.iterable, n: n, done: false}}
}

type skipIterable[T any] struct {
	iter Iterable[T]
	n    uint
	done bool
}

func (s *skipIterable[T]) Next() option.Option[T] {
	if !s.done {
		advanceByImplFiltering(s.iter, s.n)
		s.done = true
	}
	return s.iter.Next()
}

func (s *skipIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := s.iter.SizeHint()
	if lower >= s.n {
		lower -= s.n
	} else {
		lower = 0
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal >= s.n {
			upper = option.Some(upperVal - s.n)
		} else {
			upper = option.Some(uint(0))
		}
	}
	return lower, upper
}

//go:inline
func takeImpl[T any](iter Iterator[T], n uint) Iterator[T] {
	return Iterator[T]{iterable: &takeIterable[T]{iter: iter.iterable, n: n, taken: 0}}
}

type takeIterable[T any] struct {
	iter  Iterable[T]
	n     uint
	taken uint
}

func (t *takeIterable[T]) Next() option.Option[T] {
	if t.taken >= t.n {
		return option.None[T]()
	}
	item := t.iter.Next()
	if item.IsSome() {
		t.taken++
	}
	return item
}

func (t *takeIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := t.iter.SizeHint()
	if lower > t.n {
		lower = t.n
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > t.n {
			upper = option.Some(t.n)
		}
	}
	return lower, upper
}

//go:inline
func stepByImpl[T any](iter Iterator[T], step uint) Iterator[T] {
	if step == 0 {
		panic("step must be greater than 0")
	}
	return Iterator[T]{iterable: &stepByIterable[T]{iter: iter.iterable, step: step, first: true}}
}

type stepByIterable[T any] struct {
	iter  Iterable[T]
	step  uint
	first bool
}

func (s *stepByIterable[T]) Next() option.Option[T] {
	if s.first {
		s.first = false
		return s.iter.Next()
	}
	if advanceByImplFiltering(s.iter, s.step-1).IsErr() {
		return option.None[T]()
	}
	return s.iter.Next()
}

func (s *stepByIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := s.iter.SizeHint()
	if lower > 0 {
		lower = (lower + s.step - 1) / s.step
	}
	if upper.IsSome() {
		upperVal := upper.Unwrap()
		if upperVal > 0 {
			upper = option.Some((upperVal + s.step - 1) / s.step)
		}
	}
	return lower, upper
}

//go:inline
func skipWhileImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &skipWhileIterable[T]{iter: iter.iterable, predicate: predicate, done: false}}
}

type skipWhileIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
	done      bool
}

func (s *skipWhileIterable[T]) Next() option.Option[T] {
	if s.done {
		return s.iter.Next()
	}
	for {
		item := s.iter.Next()
		if item.IsNone() {
			s.done = true
			return option.None[T]()
		}
		if !s.predicate(item.Unwrap()) {
			s.done = true
			return item
		}
	}
}

func (s *skipWhileIterable[T]) SizeHint() (uint, option.Option[uint]) {
	// Can't provide accurate size hint since we don't know how many will be skipped
	return 0, option.None[uint]()
}

//go:inline
func takeWhileImpl[T any](iter Iterator[T], predicate func(T) bool) Iterator[T] {
	return Iterator[T]{iterable: &takeWhileIterable[T]{iter: iter.iterable, predicate: predicate}}
}

type takeWhileIterable[T any] struct {
	iter      Iterable[T]
	predicate func(T) bool
}

func (t *takeWhileIterable[T]) Next() option.Option[T] {
	item := t.iter.Next()
	if item.IsNone() {
		return option.None[T]()
	}
	if !t.predicate(item.Unwrap()) {
		return option.None[T]()
	}
	return item
}

func (t *takeWhileIterable[T]) SizeHint() (uint, option.Option[uint]) {
	// Can't provide accurate size hint since we don't know how many will be taken
	return 0, option.None[uint]()
}

// Skip creates an iterator that skips the first n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var skipped = iterator.Skip(2)
//	assert.Equal(t, option.Some(3), skipped.Next())
//
//go:inline
func (it Iterator[T]) Skip(n uint) Iterator[T] {
	return skipImpl(it, n)
}

// Take creates an iterator that yields the first n elements.
//
// # Examples
//
//	var iter = FromSlice([]int{1, 2, 3})
//	var taken = iterator.Take(2)
//	assert.Equal(t, option.Some(1), taken.Next())
//	assert.Equal(t, option.Some(2), taken.Next())
//	assert.Equal(t, option.None[int](), taken.Next())
//
//go:inline
func (it Iterator[T]) Take(n uint) Iterator[T] {
	return takeImpl(it, n)
}

// StepBy creates an iterator starting at the same point, but stepping by
// the given amount at each iteration.
//
// # Panics
//
// Panics if step is 0.
//
// # Examples
//
//	var iter = FromSlice([]int{0, 1, 2, 3, 4, 5})
//	var stepped = iterator.StepBy(2)
//	assert.Equal(t, option.Some(0), stepped.Next())
//	assert.Equal(t, option.Some(2), stepped.Next())
//	assert.Equal(t, option.Some(4), stepped.Next())
//
//go:inline
func (it Iterator[T]) StepBy(step uint) Iterator[T] {
	return stepByImpl(it, step)
}

// SkipWhile creates an iterator that skips elements while predicate returns true.
//
// # Examples
//
//	var iter = FromSlice([]int{-1, 0, 1})
//	var skipped = iterator.SkipWhile(func(x int) bool { return x < 0 })
//	assert.Equal(t, option.Some(0), skipped.Next())
//
//go:inline
func (it Iterator[T]) SkipWhile(predicate func(T) bool) Iterator[T] {
	return skipWhileImpl(it, predicate)
}

// TakeWhile creates an iterator that yields elements while predicate returns true.
//
// # Examples
//
//	var iter = FromSlice([]int{-1, 0, 1})
//	var taken = iterator.TakeWhile(func(x int) bool { return x < 0 })
//	assert.Equal(t, option.Some(-1), taken.Next())
//	assert.Equal(t, option.None[int](), taken.Next())
//
//go:inline
func (it Iterator[T]) TakeWhile(predicate func(T) bool) Iterator[T] {
	return takeWhileImpl(it, predicate)
}
