package errs

import (
	"errors"
	"fmt"
)

type Errs struct {
	err error
	msg string
}

func New(msg string) *Errs {
	return &Errs{
		err: errors.New(msg),
		msg: "",
	}
}

func (e *Errs) Error() string {
	if e.msg == "" {
		return e.err.Error()
	}

	return fmt.Sprintf("%s; %s", e.err.Error(), e.msg)
}

func Wrap(err error) *Errs {
	return &Errs{
		err: ErrInternal,
		msg: fmt.Sprintf("internal error: %s", err.Error()),
	}
}

func (e *Errs) Unwrap() error {
	return e.err
}
