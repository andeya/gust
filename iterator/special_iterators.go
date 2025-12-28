package iterator

import (
	"github.com/andeya/gust/option"
)

// PeekableIterator is an iterator that supports peeking at the next element.
// It embeds Iterator[T] to inherit all Iterator methods, and adds Peek() method.
//
// # Examples
//
//	var xs = []int{1, 2, 3}
//	var iter = Peekable(FromSlice(xs))
//
//	// peek() lets us see into the future
//	assert.Equal(t, option.Some(1), iterator.Peek())
//	assert.Equal(t, option.Some(1), iterator.Next())
//
//	// Can use all Iterator methods:
//	var filtered = iterator.Filter(func(x int) bool { return x > 1 })
//	assert.Equal(t, option.Some(2), filtered.Next())
type PeekableIterator[T any] struct {
	Iterator[T] // Embed Iterator to inherit all its methods
	peeker      *peekableIterable[T]
}

// Peek returns the next element without consuming it.
//
// # Examples
//
//	var xs = []int{1, 2, 3}
//	var iter = Peekable(FromSlice(xs))
//
//	assert.Equal(t, option.Some(1), iterator.Peek())
//	assert.Equal(t, option.Some(1), iterator.Peek()) // Can peek multiple times
//	assert.Equal(t, option.Some(1), iterator.Next())
//
//go:inline
func (p *PeekableIterator[T]) Peek() option.Option[T] {
	return p.peeker.Peek()
}

// Next advances the iterator and returns the next value.
// This overrides Iterator[T].Next() to handle peeked values.
//
//go:inline
func (p *PeekableIterator[T]) Next() option.Option[T] {
	return p.peeker.Next()
}

// SizeHint returns the bounds on the remaining length of the iterator.
// This overrides Iterator[T].SizeHint() to account for peeked values.
//
//go:inline
func (p *PeekableIterator[T]) SizeHint() (uint, option.Option[uint]) {
	return p.peeker.SizeHint()
}

func peekableImpl[T any](iter Iterator[T]) PeekableIterator[T] {
	core := &peekableIterable[T]{iter: iter.iterable, peeked: option.None[T]()}
	return PeekableIterator[T]{
		Iterator: Iterator[T]{iterable: core},
		peeker:   core,
	}
}

type peekableIterable[T any] struct {
	iter   Iterable[T]
	peeked option.Option[T]
}

func (p *peekableIterable[T]) Next() option.Option[T] {
	if p.peeked.IsSome() {
		item := p.peeked
		p.peeked = option.None[T]()
		return item
	}
	return p.iter.Next()
}

func (p *peekableIterable[T]) Peek() option.Option[T] {
	if p.peeked.IsNone() {
		p.peeked = p.iter.Next()
	}
	return p.peeked
}

func (p *peekableIterable[T]) SizeHint() (uint, option.Option[uint]) {
	lower, upper := p.iter.SizeHint()
	if p.peeked.IsSome() {
		if lower > 0 {
			lower++
		}
		if upper.IsSome() {
			upperVal := upper.Unwrap()
			if upperVal > 0 {
				upper = option.Some(upperVal + 1)
			}
		}
	}
	return lower, upper
}

// Cloned creates an iterator which clones all of its elements.
//
// This is useful when you have an iterator over *T, but you need an
// iterator over T.
//
// There is no guarantee whatsoever about the clone method actually
// being called *or* optimized away. So code should not depend on
// either.
//
// # Examples
//
//	var a = []string{"hello", "world"}
//	var ptrs = []*string{&a[0], &a[1]}
//	var iter = Cloned(FromSlice(ptrs))
//	var v = Collect(iter)
//	assert.Equal(t, []string{"hello", "world"}, v)
//
// clonedImpl is the internal implementation of Cloned.
//
//go:inline
func clonedImpl[T any](iter Iterable[*T]) Iterable[T] {
	return Map(Iterator[*T]{iterable: iter}, func(ptr *T) T {
		if ptr == nil {
			var zero T
			return zero
		}
		// In Go, we need to handle cloning manually for types that need it
		// For basic types, dereferencing is enough
		return *ptr
	}).iterable
}

// Cloned creates an iterator which clones all of its elements.
// This function accepts Iterator[*T] and returns Iterator[T] for chainable calls.
func Cloned[T any](iter Iterator[*T]) Iterator[T] {
	return Iterator[T]{iterable: clonedImpl(iter.iterable)}
}

//go:inline
func cycleImpl[T any](iter Iterable[T]) Iterator[T] {
	return Iterator[T]{iterable: &cycleIterable[T]{iter: iter, cache: []T{}, index: 0, exhausted: false}}
}

type cycleIterable[T any] struct {
	iter      Iterable[T]
	cache     []T
	index     int
	exhausted bool
}

func (c *cycleIterable[T]) Next() option.Option[T] {
	if !c.exhausted {
		item := c.iter.Next()
		if item.IsSome() {
			c.cache = append(c.cache, item.Unwrap())
			return item
		}
		c.exhausted = true
		if len(c.cache) == 0 {
			return option.None[T]()
		}
	}

	if len(c.cache) == 0 {
		return option.None[T]()
	}

	item := c.cache[c.index]
	c.index = (c.index + 1) % len(c.cache)
	return option.Some(item)
}

func (c *cycleIterable[T]) SizeHint() (uint, option.Option[uint]) {
	if c.exhausted {
		return 0, option.None[uint]() // Infinite
	}
	lower, upper := c.iter.SizeHint()
	if len(c.cache) > 0 {
		return 0, option.None[uint]() // Infinite
	}
	return lower, upper
}
