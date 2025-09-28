package errors

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

type StatusCode int

// Union of HTTP status codes & https://grpc.github.io/grpc/core/md_doc_statuscodes.html
// If you add a value here, search for all switches on this enum.
const (
	StatusUnknown StatusCode = iota
	StatusCanceled
	StatusInvalidArgument
	StatusDeadlineExceeded
	StatusInternal
	StatusNotFound
	StatusUnauthenticated
	StatusPermissionDenied
	StatusAlreadyExists
	StatusFailedPrecondition
	StatusUnimplemented
)

func HTTPCode(code StatusCode) int {
	switch code {
	case StatusInternal, StatusUnknown:
		return http.StatusInternalServerError
	case StatusCanceled:
		return http.StatusRequestTimeout
	case StatusUnimplemented:
		return http.StatusNotImplemented
	case StatusNotFound:
		return http.StatusNotFound
	case StatusDeadlineExceeded:
		return http.StatusGatewayTimeout
	case StatusInvalidArgument, StatusFailedPrecondition:
		return http.StatusBadRequest
	case StatusAlreadyExists:
		return http.StatusConflict
	case StatusUnauthenticated:
		return http.StatusUnauthorized
	case StatusPermissionDenied:
		return http.StatusForbidden
	default:
		log.Printf("unknown status code: %d", code)
		return http.StatusInternalServerError
	}
}

// Compile-time interface implementation checks.
var (
	_ error         = (*withCodeError)(nil)
	_ fmt.Formatter = (*withCodeError)(nil)
)

type withCodeError struct {
	code StatusCode

	// Custom message (visible to user)
	message string

	// Custom payload (visible to user)
	payload any

	cause error
}

func (e *withCodeError) Unwrap() error {
	return e.cause
}

func (e *withCodeError) Error() string {
	switch {
	case e.message != "" && e.cause != nil:
		return fmt.Sprintf("%s (%s)", e.cause.Error(), e.message)
	case e.message != "":
		return e.message
	case e.cause != nil:
		return e.cause.Error()
	default:
		return fmt.Sprintf("error %d", HTTPCode(e.code))
	}
}

func (e *withCodeError) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			_, _ = fmt.Fprintf(s, "%+v", e.Unwrap())
			if e.message != "" {
				_, _ = fmt.Fprintf(s, "\n(%s)", e.message)
			}
			return
		}
		fallthrough
	case 's':
		_, _ = io.WriteString(s, e.Error())
	case 'q':
		_, _ = fmt.Fprintf(s, "%q", e.Error())
	}
}

func WithCode(code StatusCode, message string, cause error) error {
	return &withCodeError{code, message, nil, cause}
}

func WithCodeStack(code StatusCode, message string, cause error) error {
	return &withCodeError{code, message, nil, withStackSkip(cause, 1)}
}

func Code(err error) (StatusCode, string) {
	var wc *withCodeError
	if !errors.As(err, &wc) {
		return StatusUnknown, ""
	}
	return wc.code, wc.message
}

func Status(err error) StatusCode {
	var wc *withCodeError
	if !errors.As(err, &wc) {
		return StatusUnknown
	}
	return wc.code
}

func Payload(err error) any {
	var wc *withCodeError
	if !errors.As(err, &wc) {
		return nil
	}
	return wc.payload
}

////////////////////////////////////////////////////////////
// Convenience functions
//////////////////////////////////////////////////////////.//

// Error with a code and internal message.
func ErrorCodef(code StatusCode, format string, args ...any) error {
	return WithCode(code, "", withStackSkip(fmt.Errorf(format, args...), 1))
}

// Error with a publicly visible message.
func ErrorMessagef(code StatusCode, message string, args ...any) error {
	return WithCode(code, fmt.Sprintf(message, args...), withStackSkip(fmt.Errorf(message, args...), 1))
}

func ErrorMessage(code StatusCode, message string) error {
	return WithCode(code, message, withStackSkip(errors.New(message), 1))
}

// Error with a payload.
func ErrorPayload(code StatusCode, payload any) error {
	return &withCodeError{code, "", payload, withStackSkip(errors.New(""), 1)}
}

func HTTPCodeMessage(code StatusCode, message string) (int, string) {
	httpCode := HTTPCode(code)
	if message == "" {
		message = http.StatusText(httpCode)
	}
	return httpCode, message
}
