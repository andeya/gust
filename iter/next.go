package iter

import "github.com/andeya/gust"

type Nextor[T any] interface {
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
	//
	// var a = []int{1, 2, 3};
	//
	// var iter = IterAnyFromVec(a);
	//
	// // A call to next() returns the next value...
	// assert.Equal(t, gust.Some(1), iter.Next());
	// assert.Equal(t, gust.Some(2), iter.Next());
	// assert.Equal(t, gust.Some(3), iter.Next());
	//
	// // ... and then None once it's over.
	// assert.Equal(t, gust.None[int](), iter.Next());
	//
	// // More calls may or may not return `gust.None[T]()`. Here, they always will.
	// assert.Equal(t, gust.None[int](), iter.Next());
	// assert.Equal(t, gust.None[int](), iter.Next());
	//
	Next() gust.Option[T]
}
