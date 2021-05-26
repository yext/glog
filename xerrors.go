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

// FormatStringArg is used to capture the raw format string passed to glog,
// which often contains an error message without identifying details that
// can be used to help group the message.
type FormatStringArg struct {
	Format string
}
