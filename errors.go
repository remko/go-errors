package errors

import (
	"errors"
	"fmt"
)

func Errorf(format string, args ...interface{}) error {
	return WithStack(fmt.Errorf(format, args...))
}

var New = errors.New
var As = errors.As
var Is = errors.Is
var Unwrap = errors.Unwrap
var Join = errors.Join
