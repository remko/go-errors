package errors

import (
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Bar() error {
	return WithStack(fmt.Errorf("my error message"))
}

func Foo() error {
	return Bar()
}

func Test_WithStack(t *testing.T) {
	err := Foo()

	simple := fmt.Sprintf("%v", err)
	if simple != "my error message" {
		t.Errorf("expected 'my error message', got '%s'", simple)
	}

	verbose := fmt.Sprintf("%+v", err)
	want := `my error message
MODULE/errors.Bar
	/PATH/stack_test.go:LINE
MODULE/errors.Foo
	/PATH/stack_test.go:LINE
MODULE/errors.Test_WithStack
	/PATH/stack_test.go:LINE`
	if want != cleanupStackTrace(verbose) {
		t.Errorf("expected:\n%s\ngot:\n%s", want, cleanupStackTrace(verbose))
	}
}

var fileLineRE = regexp.MustCompile(`(?m)^\t\/(.*)\/([^\/]+):(\d+)$`)
var moduleRE = regexp.MustCompile(`(?m)^[^\s].*\/([^\/]+)$`)

func cleanupStackTrace(s string) string {
	s = fileLineRE.ReplaceAllString(s, "\t/PATH/$2:LINE")
	s = moduleRE.ReplaceAllString(s, "MODULE/$1")
	return strings.Join(strings.Split(s, "\n")[0:7], "\n")
}
