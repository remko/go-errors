package errors

import (
	"errors"
	"fmt"
	"io"
	"runtime"
)

func WithStack(err error) error {
	if err == nil {
		return nil
	}
	return &withStack{
		err,
		callers(0),
	}
}

// Compile-time interface implementation checks
var (
	_ error         = (*withStack)(nil)
	_ fmt.Formatter = (*withStack)(nil)
)

type withStack struct {
	error
	*stack
}

func (w *withStack) Unwrap() error { return w.error }

func (w *withStack) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v", w.Unwrap())
			w.stack.Format(s, verb)
			return
		}
		fallthrough
	case 's':
		io.WriteString(s, w.Error())
	case 'q':
		fmt.Fprintf(s, "%q", w.Error())
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
		fmt.Fprintf(st, "\n%s\n\t%s:%d", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}
}

func HasStack(err error) bool {
	var ws withStack
	return errors.As(err, &ws)
}
