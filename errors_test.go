package errors

import (
	"errors"
	"testing"
)

var ErrNotFound1 = errors.New("not found")
var ErrForbidden1 = errors.New("forbidden")

func TestHelperErrorf(t *testing.T) {
	err := HelperErrorf("%w", ErrNotFound1)

	if !errors.Is(err, ErrNotFound1) {
		t.Errorf("expected error to be ErrNotFound1")
	}
	if errors.Is(err, ErrForbidden1) {
		t.Errorf("expected error to not be ErrForbidden1")
	}
}
