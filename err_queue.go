package xerrorz

import (
	"fmt"

	"golang.org/x/xerrors"
)

// ErrQueue is a helper struct for allowing xerrors to go along well with tree-structured errors
type ErrQueue struct {
	Curr  error // Not ErrQueue
	Queue []*error
}

func NewErrQueue(
	curr error,
	queue []*error) *ErrQueue {
	if _, isErrQueue := curr.(ErrQueue); isErrQueue {
		panic("Nested ErrQueue is not allowed") // TODO: Support ErrQueue in ErrQueue.Curr
	}

	return &ErrQueue{
		Curr:  curr,
		Queue: queue}
}

func (e ErrQueue) Error() string {
	return e.Curr.Error()
}

func (e ErrQueue) Format(s fmt.State, v rune) {
	xerrors.FormatError(e, s, v)
}

func (e ErrQueue) FormatError(p xerrors.Printer) error {
	fErr, isFormatter := e.Curr.(xerrors.Formatter)
	_, isErrQueue := e.Curr.(ErrQueue)
	if isFormatter && !isErrQueue {
		fErr.FormatError(p)
	} else {
		p.Print(e.Curr.Error())
	}

	return e.Unwrap()
}

func (e ErrQueue) Unwrap() error {
	// Dig the error if a wrapper
	if wErr, ok := e.Curr.(xerrors.Wrapper); ok {
		return ErrQueue{
			Curr:  wErr.Unwrap(),
			Queue: e.Queue}
	}

	// Dequeue
	l := len(e.Queue)
	switch {
	case l >= 2:
		return ErrQueue{
			Curr:  *e.Queue[0],
			Queue: e.Queue[1:]}
	case l == 1:
		return *e.Queue[0]
	default:
		return nil
	}
}
