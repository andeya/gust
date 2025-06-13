package iter_test

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/andeya/gust"
	"github.com/andeya/gust/digit"
	"github.com/andeya/gust/iter"
	"github.com/andeya/gust/ret"
	"github.com/stretchr/testify/assert"
)

func TestAdvanceBy(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	c <- 4
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3, 4),
		iter.FromRange(1, 4, true),
		iter.FromChan(c),
	} {
		assert.Equal(t, gust.NonErrable[uint](), i.AdvanceBy(2))
		assert.Equal(t, gust.Some(3), i.Next())
		assert.Equal(t, gust.NonErrable[uint](), i.AdvanceBy(0))
		assert.Equal(t, gust.ToErrable[uint](1), i.AdvanceBy(100)) // only `4` was skipped
	}
}

func TestAll_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var all = iter.FromVec[int8](a).All(func(v int8) bool {
		return v > 0
	})
	assert.True(t, all)
	all = iter.FromVec[int8](a).All(func(v int8) bool {
		return v > 2
	})
	assert.False(t, all)
}

func TestAll_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	var all = i.All(func(v int8) bool {
		return v != 2
	})
	assert.False(t, all)
	// we can still use `i`, as there are more elements.
	var next = i.Next()
	assert.Equal(t, gust.Some(int8(3)), next)
}

func TestAny(t *testing.T) {
	var iter = iter.FromVec([]int{1, 2, 3})
	if !iter.Any(func(x int) bool {
		return x > 1
	}) {
		t.Error("Any failed")
	}
}

func TestAny_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var hasAny = iter.FromVec[int8](a).Any(func(v int8) bool {
		return v > 0
	})
	assert.True(t, hasAny)
	hasAny = iter.FromVec[int8](a).Any(func(v int8) bool {
		return v > 5
	})
	assert.False(t, hasAny)
}

func TestAny_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	var hasAny = i.Any(func(v int8) bool {
		return v != 2
	})
	assert.True(t, hasAny)
	// we can still use `i`, as there are more elements.
	var next = i.Next()
	assert.Equal(t, gust.Some(int8(2)), next)
}

func TestChain(t *testing.T) {
	var c = make(chan int, 3)
	c <- 4
	c <- 5
	c <- 6
	close(c)
	for _, x := range [][2]iter.Iterator[int]{
		{iter.FromElements(1, 2, 3), iter.FromElements(4, 5, 6)},
		{iter.FromRange(1, 4), iter.FromChan(c)},
	} {
		var i = x[0].ToChain(x[1])
		assert.Equal(t, gust.Some(1), i.Next())
		assert.Equal(t, gust.Some(2), i.Next())
		assert.Equal(t, gust.Some(3), i.Next())
		assert.Equal(t, gust.Some(4), i.Next())
		assert.Equal(t, gust.Some(5), i.Next())
		assert.Equal(t, gust.Some(6), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestCount_1(t *testing.T) {
	assert.Equal(t, uint(3), iter.FromElements(1, 2, 3).Count())
	assert.Equal(t, uint(5), iter.FromElements(1, 2, 3, 4, 5).Count())
}

func TestCount_2(t *testing.T) {
	assert.Equal(t, uint(3), iter.FromRange(1, 3, true).Count())
	assert.Equal(t, uint(5), iter.FromRange(1, 6, false).Count())
	assert.Equal(t, uint(5), iter.FromRange(1, 6).Count())
}

func TestCount_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	assert.Equal(t, uint(3), iter.FromChan(c).Count())
}

func TestEnumerate(t *testing.T) {
	var i = iter.EnumElements[rune]('a', 'b', 'c')
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 0, Elem: 'a',
	}), i.Next())
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 1, Elem: 'b',
	}), i.Next())
	assert.Equal(t, gust.Some(gust.VecEntry[rune]{
		Index: 2, Elem: 'c',
	}), i.Next())
	assert.Equal(t, gust.None[gust.VecEntry[rune]](), i.Next())
}

func TestFilterMap_1(t *testing.T) {
	var c = make(chan string, 10)
	for idx, i := range []iter.Iterator[string]{
		iter.FromElements("1", "two", "NaN", "four", "5"),
		iter.FromChan(c),
	} {
		var i = iter.ToFilterMap[string, int](i.ToInspect(func(v string) {
			if idx == 0 {
				c <- v
			}
		}), func(v string) gust.Option[int] { return gust.Ret(strconv.Atoi(v)).Ok() })
		assert.Equal(t, gust.Some[int](1), i.Next())
		assert.Equal(t, gust.Some[int](5), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestFilterMap_2(t *testing.T) {
	var c = make(chan string, 10)
	for idx, i := range []iter.Iterator[string]{
		iter.FromElements("1", "two", "NaN", "four", "5"),
		iter.FromChan(c),
	} {
		var i = i.ToInspect(func(v string) {
			if idx == 0 {
				c <- v
			}
		}).ToXFilterMap(func(v string) gust.Option[any] { return gust.Ret(strconv.Atoi(v)).Ok().ToX() })
		assert.Equal(t, gust.Some[any](1), i.Next())
		assert.Equal(t, gust.Some[any](5), i.Next())
		assert.Equal(t, gust.None[any](), i.Next())
	}
}

func TestFilter(t *testing.T) {
	var c = make(chan int, 10)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		var i = i.ToFilter(func(v int) bool { return v > 0 })
		assert.Equal(t, gust.Some(1), i.Next())
		assert.Equal(t, gust.Some(2), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestFindMap(t *testing.T) {
	var firstNumbe = iter.FromElements("lol", "NaN", "2", "5").
		XFindMap(func(s string) gust.Option[any] {
			return gust.Ret(strconv.Atoi(s)).XOk()
		})
	assert.Equal(t, gust.Some[any](int(2)), firstNumbe)
}

func TestFind_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	assert.Equal(t, gust.Some[int8](2), i.Find(func(v int8) bool {
		return v == 2
	}))
	assert.Equal(t, gust.None[int8](), i.Find(func(v int8) bool {
		return v == 5
	}))
}

func TestFind_2(t *testing.T) {
	var a = []int8{1, 2, 3}
	var i = iter.FromVec[int8](a)
	// / assert_eq!(iter.find(|&&x| x == 2), Some(&2));
	assert.Equal(t, gust.Some[int8](2), i.Find(func(v int8) bool {
		return v == 2
	}))
	// we can still use `iter`, as there are more elements.
	assert.Equal(t, gust.Some[int8](3), i.Next())
}

func TestFlatMap(t *testing.T) {
	var c = make(chan string, 10)
	for _, i := range []iter.Iterator[string]{
		iter.FromElements("alpha", "beta", "gamma").ToInspect(func(v string) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		var merged = iter.ToFlatMap(i, func(t string) iter.Iterator[rune] {
			return iter.FromString[rune](t)
		}).Collect()
		assert.Equal(t, "alphabetagamma", string(merged))
	}
}

func TestFlatten(t *testing.T) {
	var i = iter.FromElements(
		iter.FromElements([]int{1, 2}, []int{3, 4}),
		iter.FromElements([]int{5, 6}, []int{7, 8}),
	)
	var d2 = iter.ToDeFlatten[iter.DeIterator[[]int], []int](i).Collect()
	assert.Equal(t, [][]int{{1, 2}, {3, 4}, {5, 6}, {7, 8}}, d2)
}

func TestFold_1(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 3, true),
		iter.FromChan(c),
	} {
		// the sum of all of the elements of the array
		var sum = i.Fold(int8(0), func(acc any, v int8) any {
			return acc.(int8) + v
		}).(int8)
		assert.Equal(t, int8(6), sum)
	}
}

func TestFold_2(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3, 4, 5).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 5, true),
		iter.FromChan(c),
	} {
		var zero = "0"
		var result = i.Fold(zero, func(acc any, v int8) any {
			return fmt.Sprintf("(%s + %d)", acc.(string), v)
		}).(string)
		assert.Equal(t, "(((((0 + 1) + 2) + 3) + 4) + 5)", result)
	}
}

func TestFold_3(t *testing.T) {
	var numbers = []int{1, 2, 3, 4, 5}

	// for loop:
	var result = 0
	for _, x := range numbers {
		result = result + x
	}

	// fold:
	var result2 = iter.FromVec(numbers).Fold(0, func(acc any, x int) any {
		return acc.(int) + x
	}).(int)

	// they're the same
	assert.Equal(t, result, result2)
}

func TestForEach_1(t *testing.T) {
	var c = make(chan int, 5)
	iter.FromRange(0, 5).
		ToMap(func(i int) int { return i*2 + 1 }).
		ForEach(func(i int) {
			c <- i
		})
	var v = iter.FromChan(c).Collect()
	assert.Equal(t, []int{1, 3, 5, 7, 9}, v)
}

func TestForEach_2(t *testing.T) {
	var c = make(chan int)
	go func() {
		iter.FromRange(0, 5).
			ToMap(func(i int) int { return i*2 + 1 }).
			ForEach(func(i int) {
				c <- i
			})
		close(c)
	}()
	var v = iter.FromChan(c).Collect()
	assert.Equal(t, []int{1, 3, 5, 7, 9}, v)
}

var _ gust.Iterable[int] = (*Alternate)(nil)

type Alternate struct {
	state int
}

func (a *Alternate) Next() gust.Option[int] {
	var val = a.state
	a.state = a.state + 1
	// if it's even, Some(i32), else None
	if val%2 == 0 {
		return gust.Some(val)
	}
	return gust.None[int]()
}

func TestFuse(t *testing.T) {
	var a = &Alternate{state: 0}
	var i = iter.FromIterable[int](a)
	// we can see our iterator going back and forth
	assert.Equal(t, gust.Some(0), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
	// however, once we fuse it...
	var j = i.ToFuse()
	assert.Equal(t, gust.Some(4), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())
	// it will always return `None` after the first time.
	assert.Equal(t, gust.None[int](), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())
	assert.Equal(t, gust.None[int](), j.Next())

}

func TestInspect(t *testing.T) {
	var numbers = iter.FromElements[string]("1", "2", "a", "3", "b").
		ToXMap(func(v string) any {
			return gust.Ret(strconv.Atoi(v))
		}).
		ToInspect(func(v any) {
			v.(gust.Result[int]).InspectErr(func(err error) {
				fmt.Println("Parsing error:", err)
			})
		}).
		ToFilterMap(func(v any) gust.Option[any] {
			return v.(gust.Result[int]).Ok().ToX()
		}).
		Collect()
	assert.Equal(t, []interface{}{1, 2, 3}, numbers)
}

func TestIsPartitioned(t *testing.T) {
	assert.True(t, iter.FromString[rune]("Iterator").IsPartitioned(func(r rune) bool {
		return r >= 'A' && r <= 'Z'
	}))
	assert.False(t, iter.FromString[rune]("IntoIterator").IsPartitioned(func(r rune) bool {
		return r >= 'A' && r <= 'Z'
	}))
}

func TestLast_1(t *testing.T) {
	var a = []int{1, 2, 3}
	var i = iter.FromVec(a)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestLast_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestLast_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	var i = iter.FromChan(c)
	assert.Equal(t, gust.Some(3), i.Last())
	assert.Equal(t, gust.None[int](), i.Last())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestIntersperse_1(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.ToIntersperse(100)
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.

	}
}

func TestIntersperse_2(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.ToPeekable().ToIntersperse(100)
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.
	}
}

func TestIntersperse_3(t *testing.T) {
	var c = make(chan int, 4)
	c <- 0
	c <- 1
	c <- 2
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2),
		iter.FromRange(0, 3),
		iter.FromChan(c),
	} {
		i := i.ToIntersperseWith(func() int { return 100 })
		assert.Equal(t, gust.Some(0), i.Next())     // The first element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(1), i.Next())     // The next element from `a`.
		assert.Equal(t, gust.Some(100), i.Next())   // The separator.
		assert.Equal(t, gust.Some(2), i.Next())     // The last element from `a`.
		assert.Equal(t, gust.None[int](), i.Next()) // The iterator is finished.
	}
}

func TestMap_1(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := i.ToMap(func(v int) int { return v * 2 })
		assert.Equal(t, gust.Some(2), i.Next())
		assert.Equal(t, gust.Some(4), i.Next())
		assert.Equal(t, gust.Some(6), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
	}
}

func TestMap_2(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := i.ToXMap(func(v int) any { return fmt.Sprintf("%d", v*2) })
		assert.Equal(t, gust.Some[any]("2"), i.Next())
		assert.Equal(t, gust.Some[any]("4"), i.Next())
		assert.Equal(t, gust.Some[any]("6"), i.Next())
		assert.Equal(t, gust.None[any](), i.Next())
	}
}

func TestMap_3(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		i := iter.ToMap(i, func(v int) string { return fmt.Sprintf("%d", v*2) })
		assert.Equal(t, gust.Some[string]("2"), i.Next())
		assert.Equal(t, gust.Some[string]("4"), i.Next())
		assert.Equal(t, gust.Some[string]("6"), i.Next())
		assert.Equal(t, gust.None[string](), i.Next())
	}
}

func TestMapWhile_1(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(-1, 4, 0, 1).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		i := i.ToMapWhile(func(x int) gust.Option[int] { return checkedDivide(16, x) })
		assert.Equal(t, gust.Some(-16), i.Next())
		assert.Equal(t, gust.Some(4), i.Next())
		assert.Equal(t, gust.None[int](), i.Next())
		assert.Equal(t, gust.Some(16), i.Next())
	}
}

func checkedDivide(x, y int) gust.Option[int] {
	if y == 0 {
		return gust.None[int]()
	}
	return gust.Some(x / y)
}

func TestMapWhile_2(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, -3, 4).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		a := i.ToXMapWhile(func(x int) gust.Option[any] {
			if x < 0 {
				return gust.None[any]()
			}
			return gust.Some[any](uint(x))
		}).Collect()
		assert.Equal(t, []any{uint(1), uint(2)}, a)
	}
}

func TestNextChunk(t *testing.T) {
	var iter = iter.FromVec([]int{1, 2, 3})
	assert.Equal(t, []int{1, 2}, iter.NextChunk(2).Unwrap())
	assert.Equal(t, []int{3}, iter.NextChunk(2).UnwrapErr())
	assert.Equal(t, []int{}, iter.NextChunk(2).UnwrapErr())
}

func TestNextChunk_1(t *testing.T) {
	var i = iter.FromString[rune]("中国-CN")
	assert.Equal(t, []rune{'中', '国'}, i.NextChunk(2).Unwrap())
	assert.Equal(t, []rune{'-', 'C'}, i.NextChunk(2).Unwrap())
	assert.Equal(t, []rune{'N'}, i.NextChunk(4).UnwrapErr())
	assert.Equal(t, []rune{}, i.NextChunk(1).UnwrapErr())
	assert.Equal(t, []rune{}, i.NextChunk(0).Unwrap())
}

func TestNextChunk_2(t *testing.T) {
	var i = iter.FromString[byte]("中国-CN")
	assert.Equal(t, []byte{0xe4, 0xb8, 0xad}, i.NextChunk(3).Unwrap()) // '中'
	assert.Equal(t, []byte{0xe5, 0x9b, 0xbd}, i.NextChunk(3).Unwrap()) // '国'
	assert.Equal(t, []byte{'-'}, i.NextChunk(1).Unwrap())
	assert.Equal(t, []byte{'C'}, i.NextChunk(1).Unwrap())
	assert.Equal(t, []byte{'N'}, i.NextChunk(3).UnwrapErr())
	assert.Equal(t, []byte{}, i.NextChunk(1).UnwrapErr())
	assert.Equal(t, []byte{}, i.NextChunk(0).Unwrap())
}

func TestNext_1(t *testing.T) {
	var a = []int{1, 2, 3}
	var i = iter.FromVec(a)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestNext_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestNext_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	var i = iter.FromChan(c)
	// A call to Next() returns the next value...
	assert.Equal(t, gust.Some(1), i.Next())
	assert.Equal(t, gust.Some(2), i.Next())
	assert.Equal(t, gust.Some(3), i.Next())
	// ... and then None once it's over.
	assert.Equal(t, gust.None[int](), i.Next())
	// More calls may or may not return `None`. Here, they always will.
	assert.Equal(t, gust.None[int](), i.Next())
	assert.Equal(t, gust.None[int](), i.Next())
}

func TestNth(t *testing.T) {
	var c = make(chan int, 4)
	c <- 1
	c <- 2
	c <- 3
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		assert.Equal(t, gust.Some(2), i.Nth(1))
		// Calling `Nth()` multiple times doesn't rewind the iterator:
		assert.Equal(t, gust.None[int](), i.Nth(1))
		// Returning `None` if there are less than `n + 1` elements:
		assert.Equal(t, gust.None[int](), i.Nth(10))
	}
}

func TestPartition(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		var even, odd = i.Partition(func(x int) bool { return x%2 == 0 })
		assert.Equal(t, []int{2}, even)
		assert.Equal(t, []int{1, 3}, odd)
	}
}

func TestPeekable(t *testing.T) {
	var i = iter.FromElements(1, 2, 3).ToPeekable()
	assert.Equal(t, gust.Some(1), i.Peek())
	assert.Equal(t, gust.Some(1), i.Peek())
	assert.Equal(t, gust.Some(1), i.Next())
	peeked := i.Peek()
	if peeked.IsSome() {
		p := peeked.GetOrInsertWith(nil)
		assert.Equal(t, 2, *p)
		*p = 1000
	}
	// The value reappears as the iterator continues
	assert.Equal(t, []int{1000, 3}, i.Collect())
}

func TestPosition_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	assert.Equal(t, gust.Some(1), iter.FromVec[int8](a).Position(func(v int8) bool {
		return v == 2
	}))
	assert.Equal(t, gust.None[int](), iter.FromVec[int8](a).Position(func(v int8) bool {
		return v == 5
	}))
}

func TestPosition_2(t *testing.T) {
	// Stopping at the first `true`:
	var a = iter.FromElements(1, 2, 3, 4)
	assert.Equal(t, gust.Some(1), a.Position(func(v int) bool {
		return v >= 2
	}))
	// we can still use `iter`, as there are more elements.
	assert.Equal(t, gust.Some(3), a.Next())
	// The returned index depends on iterator state
	assert.Equal(t, gust.Some(0), a.Position(func(v int) bool {
		return v == 4
	}))
}

func findMax[T gust.Integer](i iter.Iterator[T]) gust.Option[T] {
	return i.Reduce(func(acc T, v T) T {
		if acc >= v {
			return acc
		}
		return v
	})
}

func TestReduce(t *testing.T) {
	var a = []int{10, 20, 5, -23, 0}
	var b []uint
	assert.Equal(t, gust.Some(20), findMax[int](iter.FromVec(a)))
	assert.Equal(t, gust.None[uint](), findMax[uint](iter.FromVec(b)))
}

func TestRev(t *testing.T) {
	for _, i := range []iter.DeIterator[int]{
		iter.FromElements(1, 2, 3),
		iter.FromRange(1, 3, true),
	} {
		var rev = i.ToRev()
		assert.Equal(t, gust.Some(3), rev.Next())
		assert.Equal(t, gust.Some(2), rev.Next())
		assert.Equal(t, gust.Some(1), rev.Next())
		assert.Equal(t, gust.None[int](), rev.Next())
	}
}

func TestRposition_1(t *testing.T) {
	var a = []int8{1, 2, 3}
	assert.Equal(t, gust.Some(2), iter.FromVec[int8](a).Rposition(func(v int8) bool {
		return v == 3
	}))
	assert.Equal(t, gust.None[int](), iter.FromVec[int8](a).Rposition(func(v int8) bool {
		return v == 5
	}))
}

func TestRposition_2(t *testing.T) {
	var a = iter.FromElements(1, 2, 3)
	assert.Equal(t, gust.Some(1), a.Rposition(func(v int) bool {
		return v == 2
	}))
	assert.Equal(t, gust.Some(1), a.Next())
}

func TestScan(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		j := i.ToScan(1, func(state *any, x int) gust.Option[any] {
			// each iteration, we'll multiply the state by the element
			*state = (*state).(int) * x
			// then, we'll yield the negation of the state
			return gust.Some[any](-(*state).(int))
		})
		assert.Equal(t, gust.Some[any](-1), j.Next())
		assert.Equal(t, gust.Some[any](-2), j.Next())
		assert.Equal(t, gust.Some[any](-6), j.Next())
		assert.Equal(t, gust.None[any](), j.Next())
	}
}

func TestSizeHint_1(t *testing.T) {
	var i = iter.FromElements(1, 2, 3)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](2), hi)
}

func TestSizeHint_2(t *testing.T) {
	var i = iter.FromRange(1, 3, true)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](2), hi)
}

func TestSizeHint_3(t *testing.T) {
	var c = make(chan int, 3)
	c <- 1
	c <- 2
	c <- 3
	var i = iter.FromChan(c)
	var lo, hi = i.SizeHint()
	assert.Equal(t, uint(3), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
	_ = i.Next()
	lo, hi = i.SizeHint()
	assert.Equal(t, uint(2), lo)
	assert.Equal(t, gust.Some[uint](3), hi)
}

func TestSkip(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		var iter = i.ToSkip(2)
		assert.Equal(t, gust.Some(3), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}
}

func TestSkipWhile(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(-1, 0, 1).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(-1, 2),
		iter.FromChan(c),
	} {
		var iter = i.ToSkipWhile(func(v int) bool {
			return v < 0
		})
		assert.Equal(t, gust.Some(0), iter.Next())
		assert.Equal(t, gust.Some(1), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}
}

func TestStepBy(t *testing.T) {
	var c = make(chan int, 6)
	c <- 0
	c <- 1
	c <- 2
	c <- 3
	c <- 4
	c <- 5
	close(c)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(0, 1, 2, 3, 4, 5),
		iter.FromRange(0, 5, true),
		iter.FromChan(c),
	} {
		var stepBy = i.ToStepBy(2)
		assert.Equal(t, gust.Some(0), stepBy.Next())
		assert.Equal(t, gust.Some(2), stepBy.Next())
		assert.Equal(t, gust.Some(4), stepBy.Next())
		assert.Equal(t, gust.None[int](), stepBy.Next())
	}
}

func TestTake(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(1, 2, 3).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(1, 4),
		iter.FromChan(c),
	} {
		var iter = i.ToTake(2)
		assert.Equal(t, gust.Some(1), iter.Next())
		assert.Equal(t, gust.Some(2), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}
}

func TestTakeWhile(t *testing.T) {
	var c = make(chan int, 10)
	for _, i := range []iter.Iterator[int]{
		iter.FromElements(-1, 0, 1).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromRange(-1, 2),
		iter.FromChan(c),
	} {
		var iter = i.ToTakeWhile(func(v int) bool {
			return v < 0
		})
		assert.Equal(t, gust.Some(-1), iter.Next())
		assert.Equal(t, gust.None[int](), iter.Next())
	}
}

func TestTryFind(t *testing.T) {
	var a = []string{"1", "2", "lol", "NaN", "5"}
	var isMyNum = func(s string, search int) gust.Result[bool] {
		return ret.Map[int, bool](gust.Ret(strconv.Atoi(s)), func(x int) bool {
			return x == search
		})
	}
	var result = iter.FromVec[string](a).TryFind(func(s string) gust.Result[bool] {
		return isMyNum(s, 2)
	})
	assert.Equal(t, gust.Ok(gust.Some[string]("2")), result)
	result = iter.FromVec[string](a).TryFind(func(s string) gust.Result[bool] {
		return isMyNum(s, 5)
	})
	assert.True(t, result.IsErr())
}

func TestTryFold_1(t *testing.T) {
	var c = make(chan int8, 10)
	for _, i := range []iter.Iterator[int8]{
		iter.FromElements[int8](1, 2, 3).ToInspect(func(v int8) {
			c <- v
		}),
		iter.FromRange[int8](1, 4),
		iter.FromChan(c),
	} {
		// the checked sum of all the elements of the array
		var sum = i.TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
			return digit.CheckedAdd[int8](acc.(int8), v).CtrlFlow().ToX()
		})
		assert.Equal(t, gust.AnyContinue(int8(6)), sum)
	}
}

func TestTryFold_2(t *testing.T) {
	var i = iter.FromElements[int8](10, 20, 30, 100, 40, 50)
	// This sum overflows when adding the 100 element
	var sum = i.TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
		return digit.CheckedAdd[int8](acc.(int8), v).CtrlFlow().ToX()
	})
	assert.Equal(t, gust.AnyBreak(gust.Void(nil)), sum)
	// Because it short-circuited, the remaining elements are still
	// available through the iterator.
	assert.Equal(t, uint(2), i.Remaining())
	assert.Equal(t, gust.Some[int8](40), i.Next())
}

func TestTryFold_3(t *testing.T) {
	var triangular8 = iter.FromRange[int8](1, 30).TryFold(int8(0), func(acc any, v int8) gust.AnyCtrlFlow {
		return digit.CheckedAdd[int8](acc.(int8), v).XMapOrElse(
			func() any {
				return gust.AnyBreak(acc)
			}, func(sum int8) any {
				return gust.AnyContinue(sum)
			}).(gust.AnyCtrlFlow)
	})
	assert.Equal(t, gust.AnyBreak(int8(120)), triangular8)

	var triangular64 = iter.FromRange[uint64](1, 30).TryFold(uint64(0), func(acc any, v uint64) gust.AnyCtrlFlow {
		return digit.CheckedAdd(acc.(uint64), v).XMapOrElse(
			func() any {
				return gust.AnyBreak(acc)
			}, func(sum uint64) any {
				return gust.AnyContinue(sum)
			}).(gust.AnyCtrlFlow)
	})
	assert.Equal(t, gust.AnyContinue(uint64(435)), triangular64)
}

func TestTryForEach(t *testing.T) {
	var c = make(chan int, 1000)
	for _, i := range []iter.Iterator[int]{
		iter.FromRange[int](2, 100).ToInspect(func(v int) {
			c <- v
		}),
		iter.FromChan(c),
	} {
		var r = i.TryForEach(func(v int) gust.AnyCtrlFlow {
			if 323%v == 0 {
				return gust.AnyBreak(v)
			}
			return gust.AnyContinue(nil)
		})
		assert.Equal(t, gust.AnyBreak(17), r)
	}
}

func TestTryReduce_1(t *testing.T) {
	// Safely calculate the sum of a series of numbers:
	var numbers = []uint{10, 20, 5, 23, 0}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Ok(gust.Some[uint](58)), sum)
}

func TestTryReduce_2(t *testing.T) {
	// Determine when a reduction short circuited:
	var numbers = []uint{1, 2, 3, ^uint(0), 4, 5}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Err[gust.Option[uint]]("overflow"), sum)
}

func TestTryReduce_3(t *testing.T) {
	// Determine when a reduction was not performed because there are no elements:
	var numbers = []uint{}
	var sum = iter.FromVec(numbers).TryReduce(func(x, y uint) gust.Result[uint] {
		return digit.CheckedAdd(x, y).OkOr("overflow")
	})
	assert.Equal(t, gust.Ok(gust.None[uint]()), sum)
}

func TestZip_1(t *testing.T) {
	var a = iter.FromVec([]string{"x", "y", "z"})
	var b = iter.FromVec([]int{1, 2})
	var i = iter.ToZip[string, int](a, b)
	var pairs = iter.Fold[gust.Pair[string, int]](i, nil, func(acc []gust.Pair[string, int], t gust.Pair[string, int]) []gust.Pair[string, int] {
		return append(acc, t)
	})
	assert.Equal(t, []gust.Pair[string, int]{{A: "x", B: 1}, {A: "y", B: 2}}, pairs)
}

func TestZip_2(t *testing.T) {
	var c = make(chan int, 3)
	c <- 4
	c <- 5
	c <- 6
	close(c)
	for _, x := range [][2]iter.Iterator[int]{
		{iter.FromElements(1, 2, 3), iter.FromElements(4, 5, 6)},
		{iter.FromRange(1, 4), iter.FromChan(c)},
	} {
		var i = iter.ToZip(x[0], x[1])
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 1, B: 4}), i.Next())
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 2, B: 5}), i.Next())
		assert.Equal(t, gust.Some(gust.Pair[int, int]{A: 3, B: 6}), i.Next())
		assert.Equal(t, gust.None[gust.Pair[int, int]](), i.Next())
	}
}

func TestToUnique(t *testing.T) {
	var data = iter.FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{10, 20, 30, 40, 50}, iter.ToUnique[int](data).Collect())
}

func TestToDeUnique(t *testing.T) {
	var data = iter.FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{10, 20, 30, 40, 50}, iter.ToDeUnique[int](data).Collect())
	var data2 = iter.FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{50, 10, 40, 20, 30}, iter.ToDeUnique[int](data2).ToRev().Collect())
	var data3 = iter.FromElements(10, 20, 30, 20, 40, 10, 50)
	assert.Equal(t, []int{50, 10, 40, 20, 30}, iter.ToDeUnique(data3.ToRev()).Collect())
}

func TestToUniqueBy(t *testing.T) {
	var data = iter.FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"a", "bb", "ccc"}, iter.ToUniqueBy[string, int](data, func(s string) int { return len(s) }).Collect())
}

func TestToDeUniqueBy(t *testing.T) {
	var f = func(s string) int { return len(s) }
	var data = iter.FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"a", "bb", "ccc"}, iter.ToDeUniqueBy[string, int](data, f).Collect())
	var data2 = iter.FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"ccc", "c", "aa"}, iter.ToDeUniqueBy[string, int](data2, f).ToRev().Collect())
	var data3 = iter.FromElements("a", "bb", "aa", "c", "ccc")
	assert.Equal(t, []string{"ccc", "c", "aa"}, iter.ToDeUniqueBy[string, int](data3.ToRev(), f).Collect())
}
