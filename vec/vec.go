package vec

import (
	"github.com/andeya/gust"
	"github.com/andeya/gust/valconv"
)

// One try to return the first element, otherwise return zero value.
func One[T any](s []T) T {
	if len(s) > 0 {
		return s[0]
	}
	return valconv.Zero[T]()
}

// Copy creates a copy of the slice.
func Copy[T any](s []T) []T {
	if s == nil {
		return nil
	}
	r := make([]T, len(s))
	copy(r, s)
	return r
}

// CopyWithin copies part of a slice to another location in the current slice.
// @target
//
//	Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
//
// @start
//
//	Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
//
// @end
//
//	Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//	If negative, end will be counted from the end.
//	If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func CopyWithin[T any](s []T, target, start int, end ...int) {
	target = fixIndex(len(s), target, true)
	if target == len(s) {
		return
	}
	sub := Slice(s, start, end...)
	for i, v := range sub {
		s[target+i] = v
	}
}

// Every tests whether all elements in the slice pass the test implemented by the provided function.
// NOTE:
//
//	Calling this method on an empty slice will return true for any condition!
func Every[T any](s []T, fn func(k int, v T) bool) bool {
	for k, v := range s {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// Fill changes all elements in the current slice to a value, from a start index to an end index.
// @value
//
//	Zero-based index at which to copy the sequence to. If negative, target will be counted from the end.
//
// @start
//
//	Zero-based index at which to start copying elements from. If negative, start will be counted from the end.
//
// @end
//
//	Zero-based index at which to end copying elements from. CopyWithin copies up to but not including end.
//	If negative, end will be counted from the end.
//	If end is omitted, CopyWithin will copy until the last index (default to len(s)).
func Fill[T any](s []T, value T, start int, end ...int) {
	fixedStart, fixedEnd, ok := fixRange(len(s), start, end...)
	if !ok {
		return
	}
	for i := fixedStart; i < fixedEnd; i++ {
		s[i] = value
	}
}

// Filter creates a new slice with all elements that pass the test implemented by the provided function.
func Filter[T any](s []T, fn func(k int, v T) bool) []T {
	ret := make([]T, 0)
	for k, v := range s {
		if fn(k, v) {
			ret = append(ret, v)
		}
	}
	return ret
}

// Find returns the key-value of the first element in the provided slice that satisfies the provided testing function.
func Find[T any](s []T, fn func(k int, v T) bool) gust.Option[gust.VecEntry[T]] {
	for k, v := range s {
		if fn(k, v) {
			return gust.Some(gust.VecEntry[T]{Index: k, Elem: v})
		}
	}
	return gust.None[gust.VecEntry[T]]()
}

// Includes determines whether a slice includes a certain value among its entries.
// @fromIndex
//
//	The index to start the search at. Defaults to 0.
func Includes[T comparable](s []T, valueToFind T, fromIndex ...int) bool {
	return IndexOf(s, valueToFind, fromIndex...) > -1
}

// IndexOf returns the first index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//
//	The index to start the search at. Defaults to 0.
func IndexOf[T comparable](s []T, searchElement T, fromIndex ...int) int {
	idx := getFromIndex(len(s), fromIndex...)
	for k, v := range s[idx:] {
		if searchElement == v {
			return k + idx
		}
	}
	return -1
}

// LastIndexOf returns the last index at which a given element can be found in the slice, or -1 if it is not present.
// @fromIndex
//
//	The index to start the search at. Defaults to 0.
func LastIndexOf[T comparable](s []T, searchElement T, fromIndex ...int) int {
	idx := getFromIndex(len(s), fromIndex...)
	for i := len(s) - 1; i >= idx; i-- {
		if searchElement == s[i] {
			return i
		}
	}
	return -1
}

// Map creates a new slice populated with the results of calling a provided function
// on every element in the calling slice.
func Map[T any, U any](s []T, mapping func(k int, v T) U) []U {
	if s == nil {
		return nil
	}
	ret := make([]U, len(s))
	for k, v := range s {
		ret[k] = mapping(k, v)
	}
	return ret
}

// Pop removes the last element from a slice and returns that element.
// This method changes the length of the slice.
func Pop[T any](s *[]T) gust.Option[T] {
	a := *s
	if len(a) == 0 {
		return gust.None[T]()
	}
	lastIndex := len(a) - 1
	last := a[lastIndex]
	a = a[:lastIndex]
	*s = a[:len(a):len(a)]
	return gust.Some(last)
}

// Push adds one or more elements to the end of a slice and returns the new length of the slice.
func Push[T any](s *[]T, element ...T) int {
	*s = append(*s, element...)
	return len(*s)
}

// PushDistinct adds one or more new elements that do not exist in the current slice at the end.
func PushDistinct[T comparable](s []T, element ...T) []T {
L:
	for _, v := range element {
		for _, vv := range s {
			if vv == v {
				continue L
			}
		}
		s = append(s, v)
	}
	return s
}

// Reduce executes a reducer function (that you provide) on each element of the slice,
// resulting in a single output value.
// @accumulator
//
//	The accumulator accumulates callback's return values.
//	It is the accumulated value previously returned to the last invocation of the callback—or initialValue,
//	if it was supplied (see below).
//
// @initialValue
//
//	A value to use as the first argument to the first call of the callback.
//	If no initialValue is supplied, the first element in the slice will be used and skipped.
func Reduce[T any](s []T, fn func(k int, v, accumulator T) T, initialValue ...T) T {
	if len(s) == 0 {
		return valconv.Zero[T]()
	}
	start := 0
	acc := s[start]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		start += 1
	}
	for i := start; i < len(s); i++ {
		acc = fn(i, s[i], acc)
	}
	return acc
}

// ReduceRight applies a function against an accumulator and each value of the slice (from right-to-left)
// to reduce it to a single value.
// @accumulator
//
//	The accumulator accumulates callback's return values.
//	It is the accumulated value previously returned to the last invocation of the callback—or initialValue,
//	if it was supplied (see below).
//
// @initialValue
//
//	A value to use as the first argument to the first call of the callback.
//	If no initialValue is supplied, the first element in the slice will be used and skipped.
func ReduceRight[T any](s []T, fn func(k int, v, accumulator T) T, initialValue ...T) T {
	if len(s) == 0 {
		return valconv.Zero[T]()
	}
	end := len(s) - 1
	acc := s[end]
	if len(initialValue) > 0 {
		acc = initialValue[0]
	} else {
		end -= 1
	}
	for i := end; i >= 0; i-- {
		acc = fn(i, s[i], acc)
	}
	return acc
}

// Reverse reverses a slice in place.
func Reverse[T any](s []T) {
	first := 0
	last := len(s) - 1
	for first < last {
		s[first], s[last] = s[last], s[first]
		first++
		last--
	}
}

// Shift removes the first element from a slice and returns that removed element.
// This method changes the length of the slice.
func Shift[T any](s *[]T) gust.Option[T] {
	a := *s
	if len(a) == 0 {
		return gust.None[T]()
	}
	first := a[0]
	a = a[1:]
	*s = a[:len(a):len(a)]
	return gust.Some(first)
}

// Slice returns a copy of a portion of a slice into a new slice object selected
// from begin to end (end not included) where begin and end represent the index of items in that slice.
// The original slice will not be modified.
func Slice[T any](s []T, begin int, end ...int) []T {
	fixedStart, fixedEnd, ok := fixRange(len(s), begin, end...)
	if !ok {
		return []T{}
	}
	return Copy[T](s[fixedStart:fixedEnd])
}

// Some tests whether at least one element in the slice passes the test implemented by the provided function.
// NOTE:
//
//	Calling this method on an empty slice returns false for any condition!
func Some[T any](s []T, fn func(k int, v T) bool) bool {
	for k, v := range s {
		if fn(k, v) {
			return true
		}
	}
	return false
}

// Splice changes the contents of a slice by removing or replacing
// existing elements and/or adding new elements in place.
func Splice[T any](s *[]T, start, deleteCount int, items ...T) {
	a := *s
	if deleteCount < 0 {
		deleteCount = 0
	}
	start, end, _ := fixRange(len(a), start, start+1+deleteCount)
	deleteCount = end - start - 1
	for i := 0; i < len(items); i++ {
		if deleteCount > 0 {
			// replace
			a[start] = items[i]
			deleteCount--
			start++
		} else {
			// insert
			lastSlice := Copy[T](a[start:])
			items = items[i:]
			a = append(a[:start], items...)
			a = append(a[:start+len(items)], lastSlice...)
			*s = a[:len(a):len(a)]
			return
		}
	}
	if deleteCount > 0 {
		a = append(a[:start], a[start+1+deleteCount:]...)
	}
	*s = a[:len(a):len(a)]
}

// Unshift adds one or more elements to the beginning of a slice and returns the new length of the slice.
func Unshift[T any](s *[]T, element ...T) int {
	*s = append(element, *s...)
	return len(*s)
}

// UnshiftDistinct adds one or more new elements that do not exist in the current slice to the beginning
// and returns the new length of the slice.
func UnshiftDistinct[T comparable](s *[]T, element ...T) int {
	a := *s
	if len(element) == 0 {
		return len(a)
	}
	m := make(map[T]struct{}, len(element))
	r := make([]T, 0, len(a)+len(element))
L:
	for _, v := range element {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		for _, vv := range a {
			if vv == v {
				continue L
			}
		}
		r = append(r, v)
	}
	r = append(r, a...)
	*s = r[:len(r):len(r)]
	return len(r)
}

// RemoveFirst removes the first matched elements from the slice,
// and returns the new length of the slice.
func RemoveFirst[T comparable](p *[]T, elements ...T) int {
	a := *p
	m := make(map[interface{}]struct{}, len(elements))
	for _, element := range elements {
		if _, ok := m[element]; ok {
			continue
		}
		m[element] = struct{}{}
		for k, v := range a {
			if v == element {
				a = append(a[:k], a[k+1:]...)
				break
			}
		}
	}
	n := len(a)
	*p = a[:n:n]
	return n
}

// RemoveEvery removes all the elements from the slice,
// and returns the new length of the slice.
func RemoveEvery[T comparable](p *[]T, elements ...T) int {
	a := *p
	m := make(map[interface{}]struct{}, len(elements))
	for _, element := range elements {
		if _, ok := m[element]; ok {
			continue
		}
		m[element] = struct{}{}
		for i := 0; i < len(a); i++ {
			if a[i] == element {
				a = append(a[:i], a[i+1:]...)
				i--
			}
		}
	}
	n := len(a)
	*p = a[:n:n]
	return n
}

// Concat is used to merge two or more slices.
// This method does not change the existing slices, but instead returns a new slice.
func Concat[T any](s ...[]T) []T {
	var totalLen int
	for _, v := range s {
		totalLen += len(v)
	}
	ret := make([]T, totalLen)
	dst := ret
	for _, v := range s {
		n := copy(dst, v)
		dst = dst[n:]
	}
	return ret
}

// Intersect calculates intersection of two or more slices,
// and returns the count of each element.
func Intersect[T comparable](s ...[]T) (intersectCount map[T]int) {
	if len(s) == 0 {
		return nil
	}
	for _, v := range s {
		if len(v) == 0 {
			return nil
		}
	}
	counts := make([]map[T]int, len(s))
	for k, v := range s {
		counts[k] = vecDistinct(v, nil)
	}
	intersectCount = counts[0]
L:
	for k, v := range intersectCount {
		for _, c := range counts[1:] {
			v2 := c[k]
			if v2 == 0 {
				delete(intersectCount, k)
				continue L
			}
			if v > v2 {
				v = v2
			}
		}
		intersectCount[k] = v
	}
	return intersectCount
}

// Distinct calculates the count of each different element,
// and only saves these different elements in place if changeSlice is true.
func Distinct[T comparable](s *[]T, changeSlice bool) (distinctCount map[T]int) {
	if !changeSlice {
		return vecDistinct(*s, nil)
	}
	a := (*s)[:0]
	distinctCount = vecDistinct(*s, &a)
	n := len(distinctCount)
	*s = a[:n:n]
	return distinctCount
}

func vecDistinct[T comparable](src []T, dst *[]T) map[T]int {
	m := make(map[T]int, len(src))
	if dst == nil {
		for _, v := range src {
			n := m[v]
			m[v] = n + 1
		}
	} else {
		a := *dst
		for _, v := range src {
			n := m[v]
			m[v] = n + 1
			if n == 0 {
				a = append(a, v)
			}
		}
		*dst = a
	}
	return m
}

// DistinctBy deduplication in place according to the mapping function
func DistinctBy[T comparable, U comparable](s *[]T, mapping func(k int, v T) U) {
	a := (*s)[:0]
	distinctCount := vecDistinctBy(*s, &a, mapping)
	n := len(distinctCount)
	*s = a[:n:n]
}

func vecDistinctBy[T comparable, U comparable](src []T, dst *[]T, mapping func(k int, v T) U) map[U]int {
	m := make(map[U]int, len(src))
	if dst == nil {
		for k, v := range src {
			x := mapping(k, v)
			n := m[x]
			m[x] = n + 1
		}
	} else {
		a := *dst
		for k, v := range src {
			x := mapping(k, v)
			n := m[x]
			m[x] = n + 1
			if n == 0 {
				a = append(a, v)
			}
		}
		*dst = a
	}
	return m
}

// DistinctMap returns the unique elements after mapping.
func DistinctMap[T any, U comparable](s []T, mapping func(k int, v T) U) []U {
	if s == nil {
		return nil
	}
	m := make(map[U]struct{}, len(s))
	a := make([]U, 0, len(m))
	for k, v := range s {
		x := mapping(k, v)
		if _, ok := m[x]; ok {
			continue
		}
		a = append(a, x)
		m[x] = struct{}{}
	}
	return a
}

// SetsUnion calculates between multiple collections: set1 ∪ set2 ∪ others...
// This method does not change the existing slices, but instead returns a new slice.
func SetsUnion[T comparable](set1, set2 []T, others ...[]T) []T {
	m := make(map[T]struct{}, len(set1)+len(set2))
	r := make([]T, 0, len(m))
	for _, set := range append([][]T{set1, set2}, others...) {
		for _, v := range set {
			_, ok := m[v]
			if ok {
				continue
			}
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}

// SetsIntersect calculates between multiple collections: set1 ∩ set2 ∩ others...
// This method does not change the existing slices, but instead returns a new slice.
func SetsIntersect[T comparable](set1, set2 []T, others ...[]T) []T {
	sets := append([][]T{set2}, others...)
	setsCount := make([]map[T]int, len(sets))
	for k, v := range sets {
		setsCount[k] = vecDistinct(v, nil)
	}
	m := make(map[T]struct{}, len(set1))
	r := make([]T, 0, len(m))
L:
	for _, v := range set1 {
		if _, ok := m[v]; ok {
			continue
		}
		m[v] = struct{}{}
		for _, m2 := range setsCount {
			if m2[v] == 0 {
				continue L
			}
		}
		r = append(r, v)
	}
	return r
}

// SetsDifference calculates between multiple collections: set1 - set2 - others...
// This method does not change the existing slices, but instead returns a new slice.
func SetsDifference[T comparable](set1, set2 []T, others ...[]T) []T {
	m := make(map[T]struct{}, len(set1))
	r := make([]T, 0, len(set1))
	sets := append([][]T{set2}, others...)
	for _, v := range sets {
		inter := SetsIntersect(set1, v)
		for _, vv := range inter {
			m[vv] = struct{}{}
		}
	}
	for _, v := range set1 {
		if _, ok := m[v]; !ok {
			r = append(r, v)
			m[v] = struct{}{}
		}
	}
	return r
}
