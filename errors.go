package errors

import (
	"errors"
	"fmt"
)

func Errorf(format string, args ...interface{}) error {
	return WithStack(fmt.Errorf(format, args...))
}

func HelperErrorf(format string, args ...any) error {
	return withStackSkip(fmt.Errorf(format, args...), 2)
}

var New = errors.New
var As = errors.As
var Is = errors.Is
var Unwrap = errors.Unwrap
var Join = errors.Join
