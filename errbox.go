package gust

import (
	"errors"
	"fmt"
)

var _ error = (*ErrBox)(nil)

// ToErrBox wraps any error type.
func ToErrBox(val any) *ErrBox {
	return &ErrBox{val: val}
}

func toError(e any) error {
	if x, is := e.(error); is {
		return x
	}
	return ToErrBox(e)
}

// ErrBox is a wrapper for any error type.
type ErrBox struct {
	val any
}

// Value returns the inner value.
func (a *ErrBox) Value() any {
	if a == nil {
		return nil
	}
	return a.val
}

// Error returns the string representation.
func (a *ErrBox) Error() string {
	if a == nil {
		return ""
	}
	switch val := a.val.(type) {
	case string:
		return val
	case error:
		return val.Error()
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", val)
	}
}

// Unwrap returns the inner error.
func (a *ErrBox) Unwrap() error {
	if a == nil {
		return nil
	}
	switch val := a.val.(type) {
	case nil:
		return nil
	case error:
		u, ok := val.(interface {
			Unwrap() error
		})
		if ok {
			return u.Unwrap()
		}
		return val
	default:
		return nil
	}
}

func (a *ErrBox) Is(target error) bool {
	switch t := target.(type) {
	case *ErrBox:
		if t == nil {
			return a == nil
		}
		if a == nil {
			return false
		}
		return a.val == t.val
	default:
		b := a.Unwrap()
		if b != nil {
			return errors.Is(b, target)
		}
		return false
	}
}

func (a *ErrBox) As(target any) bool {
	if target == nil {
		panic("errors: target cannot be nil")
	}
	t, ok := target.(*ErrBox)
	if ok {
		if t != nil {
			t.val = a.val
			return true
		}
		return false
	}
	b := a.Unwrap()
	if b != nil {
		return errors.As(b, target)
	}
	return false
}
