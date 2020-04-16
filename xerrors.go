package glog

import "golang.org/x/xerrors"

// ErrorArg captures information about an error passed as an argument.
// It is passed to backends as a data arg.
type ErrorArg struct {
	Error error
}

// RootCause returns the innermost error.
func (xe ErrorArg) RootCause() error {
	err := xe.Error
	for {
		wrapper, ok := err.(xerrors.Wrapper)
		if !ok {
			break
		}
		err = wrapper.Unwrap()
	}
	return err
}

// InterestingCause returns the error with the most succinct, yet interesting
// message. The RootCause may be something boring like "context canceled". The
// error returned here strives for something better.
func (xe ErrorArg) InterestingCause() error {
	var prev error
	var err = xe.Error
	for {
		wrapper, ok := err.(xerrors.Wrapper)
		if !ok {
			break
		}
		prev = err
		err = wrapper.Unwrap()
	}
	if prev != nil {
		return prev
	}
	return err
}
