package errors

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

func WithStack(err error) error {
	return withStackSkip(err, 1)
}

func withStackSkip(err error, skip int) error {
	if err == nil {
		return nil
	}
	return &withStackError{
		err,
		callers(skip),
	}
}

// Compile-time interface implementation checks.
var (
	_ error         = (*withStackError)(nil)
	_ fmt.Formatter = (*withStackError)(nil)
)

type withStackError struct {
	error
	*stack
}

func (w *withStackError) Unwrap() error { return w.error }

func (w *withStackError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, w.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", w.Error())
	}
}

func callers(skip int) *stack {
	const depth = 32
	var pcs [depth]uintptr
	n := runtime.Callers(skip+3, pcs[:])
	var st stack = pcs[:n]
	return &st
}

type stack []uintptr

func (s *stack) Format(st fmt.State, verb rune) {
	frames := runtime.CallersFrames(*s)
	for {
		frame, more := frames.Next()
		_, _ = fmt.Fprintf(st, "\n%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func HasStack(err error) bool {
	var ws withStackError
	return errors.As(err, &ws)
}
