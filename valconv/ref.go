package valconv

// Zero returns the zero value of the type.
func Zero[T any]() T {
	var t T
	return t
}

// Ref returns the address of the value.
func Ref[T any](t T) *T {
	return &t
}

// Deref returns the value of the pointer type.
// NOTE:
//
//	If the value is nil, it will return the zero T.
func Deref[T any](t *T) T {
	if t == nil {
		var x T
		return x
	}
	return *t
}

// RefSlice convert []T to []*T.
func RefSlice[T any](a []T) []*T {
	if a == nil {
		return nil
	}
	s := make([]*T, len(a))
	for i, t := range a {
		s[i] = &t
	}
	return s
}

// DerefSlice convert []*T to []T.
// NOTE:
//
// If an element is nil, it will be set to the zero T.
func DerefSlice[T any](a []*T) []T {
	if a == nil {
		return nil
	}
	s := make([]T, len(a))
	for i, t := range a {
		if t != nil {
			s[i] = *t
		}
	}
	return s
}

// RefMap convert map[K]V to map[K]*V.
func RefMap[K comparable, V any](a map[K]V) map[K]*V {
	if a == nil {
		return nil
	}
	s := make(map[K]*V, len(a))
	for k, v := range a {
		s[k] = &v
	}
	return s
}

// DerefMap convert map[K]*V to map[K]V.
// NOTE:
//
// If a value is nil, it will be set to the zero V.
func DerefMap[K comparable, V any](a map[K]*V) map[K]V {
	if a == nil {
		return nil
	}
	s := make(map[K]V, len(a))
	for k, v := range a {
		if v != nil {
			s[k] = *v
		} else {
			var x V
			s[k] = x
		}
	}
	return s
}
