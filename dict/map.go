// Package dict is a package of generic-type functions for map.
package dict

import "github.com/andeya/gust/option"

// DictEntry is a key-value entry of map.
type DictEntry[K comparable, V any] struct {
	Key   K
	Value V
}

// Split splits the dictionary entry into its two components.
//
//go:inline
func (d DictEntry[K, V]) Split() (K, V) {
	return d.Key, d.Value
}

// Get returns the option.Option[V] of the entry for the provided key.
func Get[K comparable, V any](m map[K]V, k K) option.Option[V] {
	if v, ok := m[k]; ok {
		return option.Some[V](v)
	}
	return option.None[V]()
}

// Keys returns the keys of map.
func Keys[K comparable, V any](m map[K]V) []K {
	if m == nil {
		return nil
	}
	ret := make([]K, 0, len(m))
	for k := range m {
		ret = append(ret, k)
	}
	return ret
}

// Values returns the values of map.
func Values[K comparable, V any](m map[K]V) []V {
	if m == nil {
		return nil
	}
	ret := make([]V, 0, len(m))
	for _, v := range m {
		ret = append(ret, v)
	}
	return ret
}

// Entries returns the entries of map.
func Entries[K comparable, V any](m map[K]V) []DictEntry[K, V] {
	if m == nil {
		return nil
	}
	ret := make([]DictEntry[K, V], 0, len(m))
	for k, v := range m {
		ret = append(ret, DictEntry[K, V]{Key: k, Value: v})
	}
	return ret
}

// Vec generates an orderless slice through the set function.
func Vec[K comparable, V any, T any](m map[K]V, set func(K, V) T) []T {
	if m == nil {
		return nil
	}
	ret := make([]T, 0, len(m))
	for k, v := range m {
		ret = append(ret, set(k, v))
	}
	return ret
}

// Copy creates a copy of the map.
func Copy[K comparable, V any](m map[K]V) map[K]V {
	if m == nil {
		return nil
	}
	r := make(map[K]V, len(m))
	for k, v := range m {
		r[k] = v
	}
	return r
}

// Every tests whether all entries in the map pass the test implemented by the provided function.
// NOTE:
//
//	Calling this method on an empty map will return true for any condition!
func Every[K comparable, V any](m map[K]V, fn func(k K, v V) bool) bool {
	for k, v := range m {
		if !fn(k, v) {
			return false
		}
	}
	return true
}

// Some tests whether at least one entry in the map passes the test implemented by the provided function.
// NOTE:
//
//	Calling this method on an empty map returns false for any condition!
func Some[K comparable, V any](m map[K]V, fn func(K, V) bool) bool {
	for k, v := range m {
		if fn(k, v) {
			return true
		}
	}
	return false
}

// Find returns an entry in the provided map that satisfies the test function.
func Find[K comparable, V any](m map[K]V, fn func(K, V) bool) option.Option[DictEntry[K, V]] {
	for k, v := range m {
		if fn(k, v) {
			return option.Some(DictEntry[K, V]{Key: k, Value: v})
		}
	}
	return option.None[DictEntry[K, V]]()
}

// Filter creates a new map with all elements that pass the test implemented by the provided function.
func Filter[K comparable, V any](m map[K]V, fn func(K, V) bool) map[K]V {
	ret := make(map[K]V, 0)
	for k, v := range m {
		if fn(k, v) {
			ret[k] = v
		}
	}
	return ret
}

// FilterMap returns a filtered and mapped map of new entries.
func FilterMap[K comparable, V any, K2 comparable, V2 any](m map[K]V, fn func(K, V) option.Option[DictEntry[K2, V2]]) map[K2]V2 {
	ret := make(map[K2]V2, 0)
	for k, v := range m {
		fn(k, v).Inspect(func(p DictEntry[K2, V2]) {
			ret[p.Key] = p.Value
		})
	}
	return ret
}

// FilterMapKey returns a filtered and mapped map of new entries.
func FilterMapKey[K comparable, V any, K2 comparable](m map[K]V, fn func(K, V) option.Option[DictEntry[K2, V]]) map[K2]V {
	return FilterMap[K, V, K2, V](m, fn)
}

// FilterMapValue returns a filtered and mapped map of new entries.
func FilterMapValue[K comparable, V any, V2 any](m map[K]V, fn func(K, V) option.Option[DictEntry[K, V2]]) map[K]V2 {
	return FilterMap[K, V, K, V2](m, fn)
}

// Map creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func Map[K comparable, V any, K2 comparable, V2 any](m map[K]V, mapping func(K, V) DictEntry[K2, V2]) map[K2]V2 {
	if m == nil {
		return nil
	}
	ret := make(map[K2]V2, len(m))
	for k, v := range m {
		x := mapping(k, v)
		ret[x.Key] = x.Value
	}
	return ret
}

// MapCurry creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func MapCurry[K comparable, V any, K2 comparable, V2 any](m map[K]V, keyMapping func(K) K2) func(valueMapping func(V) V2) map[K2]V2 {
	return func(valueMapping func(V) V2) map[K2]V2 {
		return Map[K, V, K2, V2](m, func(k K, v V) DictEntry[K2, V2] {
			return DictEntry[K2, V2]{
				Key:   keyMapping(k),
				Value: valueMapping(v),
			}
		})
	}
}

// MapKey creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func MapKey[K comparable, V any, K2 comparable](m map[K]V, mapping func(K, V) K2) map[K2]V {
	return Map[K, V, K2, V](m, func(k K, v V) DictEntry[K2, V] {
		return DictEntry[K2, V]{
			Key:   mapping(k, v),
			Value: v,
		}
	})
}

// MapKeyAlone creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func MapKeyAlone[K comparable, V any, K2 comparable](m map[K]V, mapping func(K) K2) map[K2]V {
	return Map[K, V, K2, V](m, func(k K, v V) DictEntry[K2, V] {
		return DictEntry[K2, V]{
			Key:   mapping(k),
			Value: v,
		}
	})
}

// MapValue creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func MapValue[K comparable, V any, V2 any](m map[K]V, mapping func(K, V) V2) map[K]V2 {
	return Map[K, V, K, V2](m, func(k K, v V) DictEntry[K, V2] {
		return DictEntry[K, V2]{
			Key:   k,
			Value: mapping(k, v),
		}
	})
}

// MapValueAlone creates a new map populated with the results of calling a provided function
// on every entry in the calling map.
func MapValueAlone[K comparable, V any, V2 any](m map[K]V, mapping func(V) V2) map[K]V2 {
	return Map[K, V, K, V2](m, func(k K, v V) DictEntry[K, V2] {
		return DictEntry[K, V2]{
			Key:   k,
			Value: mapping(v),
		}
	})
}
