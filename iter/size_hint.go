package iter

import "github.com/andeya/gust"

type SizeHint interface {
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
	//
	// var a = []int{1, 2, 3};
	// var iter = IterAnyFromVec(a);
	//
	// assert.Equal(t, (3, gust.Some(3)), iter.SizeHint());
	//
	//
	// A more complex example:
	//
	//
	// // The even numbers in the range of zero to nine.
	// var iter = IterAnyFromRange(0..10).Filter(func(x T) {return x % 2 == 0});
	//
	// // We might iterate from zero to ten times. Knowing that it's five
	// // exactly wouldn't be possible without executing filter().
	// assert.Equal(t, (0, gust.Some(10)), iter.SizeHint());
	//
	// // Let's add five more numbers with chain()
	// var iter = IterAnyFromRange(0, 10).Filter(func(x T) {return x % 2 == 0}).Chain(IterAnyFromRange(15, 20));
	//
	// // now both bounds are increased by five
	// assert.Equal(t, (5, gust.Some(15)), iter.SizeHint());
	//
	//
	// Returning `gust.None[int]()` for an upper bound:
	//
	//
	// // an infinite next has no upper bound
	// // and the maximum possible lower bound
	// var iter = IterAnyFromRange(0, math.MaxInt);
	//
	// assert.Equal(t, (math.MaxInt, gust.None[int]()), iter.SizeHint());
	//
	SizeHint() (uint64, gust.Option[uint64])
}
