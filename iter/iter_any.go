package iter

import (
	"github.com/andeya/gust"
)

func newAnyIter[T any](step uint, filterMap func(T) gust.Option[T], nextList ...Nextor[T]) *AnyIter[T] {
	iter := &AnyIter[T]{nextChain: nextList, nextChainIndex: 0, filterMap: filterMap, firstTake: true, step: step}
	iter.iterTrait = iterTrait[T, iterCore[T]]{
		core: iter,
	}
	return iter
}

type (
	AnyIter[T any] struct {
		iterTrait[T, iterCore[T]]
		step           uint
		nextChainIndex int
		nextChain      []Nextor[T]
		filterMap      func(T) gust.Option[T]
		firstTake      bool
	}
	counter interface {
		count() uint64
	}
)

// Next advances the next and returns the next value.
//
// Returns [`gust.None[T]()`] when iteration is finished. Individual next
// implementations may choose to resume iteration, and so calling `next()`
// again may or may not eventually min returning [`gust.Some(T)`] again at some
// point.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
//
// var iter = IterAnyFromVec(a);
//
// A call to next() returns the next value...
// assert.Equal(t, gust.Some(1), iter.Next());
// assert.Equal(t, gust.Some(2), iter.Next());
// assert.Equal(t, gust.Some(3), iter.Next());
//
// ... and then None once it's over.
// assert.Equal(t, gust.None[int](), iter.Next());
//
// More calls may or may not return `gust.None[T]()`. Here, they always will.
// assert.Equal(t, gust.None[int](), iter.Next());
// assert.Equal(t, gust.None[int](), iter.Next());
func (iter *AnyIter[T]) Next() gust.Option[T] {
	if iter.firstTake {
		iter.firstTake = false
		return iter.elemNext()
	} else {
		switch iter.step {
		case 0:
			return gust.None[T]()
		case 1:
			return iter.elemNext()
		default:
			var v gust.Option[T]
			for i := iter.step; i > 0; i-- {
				v = iter.elemNext()
			}
			return v
		}
	}
}

func (iter *AnyIter[T]) elemNext() gust.Option[T] {
	if iter.nextChainIndex < len(iter.nextChain) {
		for iter.nextChainIndex < len(iter.nextChain) {
			next := iter.nextChain[iter.nextChainIndex]
			v := next.Next()
			if iter.filterMap != nil {
				for v.IsSome() {
					v = iter.filterMap(v.Unwrap())
					if v.IsSome() {
						break
					}
					v = next.Next()
				}
			}
			if v.IsSome() {
				return v
			}
			iter.nextChainIndex++
		}
	}
	return gust.None[T]()
}

// SizeHint returns the bounds on the remaining length of the next.
//
// Specifically, `SizeHint()` returns a tuple where the first element
// is the lower bound, and the second element is the upper bound.
//
// The second half of the tuple that is returned is an <code>Option[T]</code>.
// A [`gust.None[T]()`] here means that either there is no known upper bound, or the
// upper bound is larger than [`int`].
//
// # Implementation notes
//
// It is not enforced that a next implementation yields the declared
// number of elements. A buggy next may yield less than the lower bound
// or more than the upper bound of elements.
//
// `SizeHint()` is primarily intended to be used for optimizations such as
// reserving space for the elements of the next, but must not be
// trusted to e.g., omit bounds checks in unsafe code. An incorrect
// implementation of `SizeHint()` should not lead to memory safety
// violations.
//
// That said, the implementation should provide a correct estimation,
// because otherwise it would be a violation of the interface's protocol.
//
// The default implementation returns <code>(0, [None[int]()])</code> which is correct for any
// next.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
// var iter = IterAnyFromVec(a);
//
// assert.Equal(t, (3, gust.Some(3)), iter.SizeHint());
//
// A more complex example:
//
// The even numbers in the range of zero to nine.
// var iter = IterAnyFromRange(0..10).Filter(func(x T) {return x % 2 == 0});
//
// We might iterate from zero to ten times. Knowing that it's five
// exactly wouldn't be possible without executing filter().
// assert.Equal(t, (0, gust.Some(10)), iter.SizeHint());
//
// Let's add five more numbers with chain()
// var iter = IterAnyFromRange(0, 10).Filter(func(x T) {return x % 2 == 0}).Chain(IterAnyFromRange(15, 20));
//
// now both bounds are increased by five
// assert.Equal(t, (5, gust.Some(15)), iter.SizeHint());
//
// Returning `gust.None[int]()` for an upper bound:
//
// an infinite next has no upper bound
// and the maximum possible lower bound
// var iter = IterAnyFromRange(0, math.MaxInt);
//
// assert.Equal(t, (math.MaxInt, gust.None[int]()), iter.SizeHint());
func (iter *AnyIter[T]) SizeHint() (uint64, gust.Option[uint64]) {
	if iter.nextChainIndex >= len(iter.nextChain) {
		return 0, gust.None[uint64]()
	}
	var a uint64
	var b uint64
	var none = false
	for _, next := range iter.nextChain[iter.nextChainIndex:] {
		if sizeHint, ok := next.(SizeHint); ok {
			x, y := sizeHint.SizeHint()
			a += x
			if none {
				continue
			}
			if y.IsNone() {
				none = true
			} else {
				b += y.Unwrap()
			}
		}
	}
	if none {
		return a, gust.None[uint64]()
	}
	return a, gust.Some(b)
}

// Count consumes the next, counting the number of iterations and returning it.
//
// This method will call [`Nextor`] repeatedly until [`gust.None[T]()`] is encountered,
// returning the number of times it saw [`gust.Some`]. Note that [`Nextor`] has to be
// called at least once even if the next does not have any elements.
//
// # Overflow Behavior
//
// The method does no guarding against overflows, so counting elements of
// a next with more than [`math.MaxInt`] elements either produces the
// wrong result or panics. If debug assertions are enabled, a panic is
// guaranteed.
//
// # Panics
//
// This function might panic if the next has more than [`math.MaxInt`]
// elements.
//
// # Examples
//
// Basic usage:
//
// var a = []int{1, 2, 3};
// assert.Equal(t, IterAnyFromVec(a).Count(), 3);
//
// var a = []int{1, 2, 3, 4, 5};
// assert.Equal(t, IterAnyFromVec(a).Count(), 5);
func (iter *AnyIter[T]) Count() uint64 {
	if iter.nextChainIndex >= len(iter.nextChain) {
		return 0
	}
	var a uint64
	for _, next := range iter.nextChain[iter.nextChainIndex:] {
		if c, ok := next.(counter); ok {
			a += c.count()
		} else {
			for next.Next().IsSome() {
				a++
			}
		}
	}
	iter.nextChainIndex = len(iter.nextChain)
	return a
}

// StepBy creates a next starting at the same point, but stepping by
// the given amount at each iteration.
//
// Note 1: The first element of the next will always be returned,
// regardless of the step given.
//
// Note 2: The time at which ignored elements are pulled is not fixed.
// `StepBy` behaves like the sequence `iter.Next()`, `iter.Nth(step-1)`,
// `iter.Nth(step-1)`, …, but is also free to behave like the sequence.
//
// # Examples
//
// Basic usage:
//
// var a = []int{0, 1, 2, 3, 4, 5};
// var iter = IterAnyVec(a).StepBy(2);
//
// assert.Equal(t, iter.Next(), Some(0));
// assert.Equal(t, iter.Next(), Some(2));
// assert.Equal(t, iter.Next(), Some(4));
// assert.Equal(t, iter.Next(), None[T]());
func (iter *AnyIter[T]) StepBy(step uint) *AnyIter[T] {
	return newAnyIter[T](step, nil, iter)
}

func (iter *AnyIter[T]) Append(other Nextor[T]) {
	iter.nextChain = append(iter.nextChain, other)
}

// Filter creates an iterator which uses a closure to determine if an element
// should be yielded.
//
// Given an element the closure must return `true` or `false`. The returned
// iterator will yield only the elements for which the closure returns
// true.
//
// # Examples
//
// Basic usage:
//
// ```
// var a = []int{0, 1, 2};
//
// var iter = IterAnyFromVec(a).Filter(func(x int)bool{return x>0});
//
// assert_eq!(iter.Next(), gust.Some(&1));
// assert_eq!(iter.Next(), gust.Some(&2));
// assert_eq!(iter.Next(), gust.None[int]());
// ```
//
// Note that `iter.Filter(f).Next()` is equivalent to `iter.Find(f)`.
func (iter *AnyIter[T]) Filter(f func(T) bool) *AnyIter[T] {
	return newAnyIter[T](1, func(t T) gust.Option[T] {
		if f(t) {
			return gust.Some(t)
		}
		return gust.None[T]()
	}, iter)
}

func (iter *AnyIter[T]) FilterMap(f func(T) gust.Option[T]) *AnyIter[T] {
	return newAnyIter[T](1, f, iter)
}
