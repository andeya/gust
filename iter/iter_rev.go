// nolint:unused
package iter

import (
	"github.com/andeya/gust"
)

var (
	_ DeIterator[any]         = (*deRevIterator[any])(nil)
	_ iRealNext[any]          = (*deRevIterator[any])(nil)
	_ iRealSizeHint           = (*deRevIterator[any])(nil)
	_ iRealAdvanceBy[any]     = (*deRevIterator[any])(nil)
	_ iRealNth[any]           = (*deRevIterator[any])(nil)
	_ iRealTryFold[any]       = (*deRevIterator[any])(nil)
	_ iRealFold[any]          = (*deRevIterator[any])(nil)
	_ iRealFind[any]          = (*deRevIterator[any])(nil)
	_ iRealRemaining          = (*deRevIterator[any])(nil)
	_ iRealNextBack[any]      = (*deRevIterator[any])(nil)
	_ iRealAdvanceBackBy[any] = (*deRevIterator[any])(nil)
	_ iRealNthBack[any]       = (*deRevIterator[any])(nil)
	_ iRealTryRfold[any]      = (*deRevIterator[any])(nil)
	_ iRealRfold[any]         = (*deRevIterator[any])(nil)
	_ iRealRfind[any]         = (*deRevIterator[any])(nil)
)

func newDeRevIterator[T any](iter DeIterator[T]) DeIterator[T] {
	p := &deRevIterator[T]{}
	p.iter = iter
	p.setFacade(p)
	return p
}

type deRevIterator[T any] struct {
	deIterBackground[T]
	iter DeIterator[T]
}

func (d *deRevIterator[T]) realNext() gust.Option[T] {
	return d.iter.NextBack()
}

func (d *deRevIterator[T]) realSizeHint() (uint, gust.Option[uint]) {
	return d.iter.SizeHint()
}

func (d *deRevIterator[T]) realAdvanceBy(n uint) gust.Errable[uint] {
	return d.iter.AdvanceBackBy(n)
}

func (d *deRevIterator[T]) realNth(n uint) gust.Option[T] {
	return d.iter.NthBack(n)
}

func (d *deRevIterator[T]) realTryFold(acc any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return d.iter.TryRfold(acc, fold)
}

func (d *deRevIterator[T]) realFold(acc any, fold func(any, T) any) any {
	return d.iter.Rfold(acc, fold)
}

func (d *deRevIterator[T]) realFind(predicate func(T) bool) gust.Option[T] {
	return d.iter.Rfind(predicate)
}

func (d *deRevIterator[T]) realRemaining() uint {
	return d.iter.Remaining()
}

func (d *deRevIterator[T]) realNextBack() gust.Option[T] {
	return d.iter.Next()
}

func (d *deRevIterator[T]) realAdvanceBackBy(n uint) gust.Errable[uint] {
	return d.iter.AdvanceBy(n)
}

func (d *deRevIterator[T]) realNthBack(n uint) gust.Option[T] {
	return d.iter.Nth(n)
}

func (d *deRevIterator[T]) realTryRfold(acc any, fold func(any, T) gust.AnyCtrlFlow) gust.AnyCtrlFlow {
	return d.iter.TryFold(acc, fold)
}

func (d *deRevIterator[T]) realRfold(acc any, fold func(any, T) any) any {
	return d.iter.Fold(acc, fold)
}

func (d *deRevIterator[T]) realRfind(predicate func(T) bool) gust.Option[T] {
	return d.iter.Find(predicate)
}
