package gust

import "fmt"

type errorWithVal struct {
	val any
}

func newAnyError(val any) error {
	if err, _ := val.(error); err != nil {
		return err
	}
	return &errorWithVal{val: val}
}

func (a *errorWithVal) Error() string {
	return fmt.Sprintf("%v", a.val)
}
