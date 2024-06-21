package errors

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	ErrUnknown        = Error("err_unknown: Unknown error occurred")
	ErrBadRequest     = Error("err_bad_request: Bad request")
	ErrNotFound       = Error("err_not_found: Not found")
	ErrInternalServer = Error("err_internal_server: Internal server error")
	ErrValidation     = Error("err_validation: Validation failed")
)

const ErrSeparator = " -- "

type Error string

func (e Error) Error() string {
	return string(e)
}

func (e Error) Is(target error) bool {
	return e.Error() == target.Error() || strings.HasPrefix(target.Error(), e.Error()+ErrSeparator)
}

func (e Error) As(target any) bool {
	v := reflect.ValueOf(target).Elem()

	if v.Type().Name() == "Error" && v.CanSet() {
		v.SetString(string(e))
		return true
	}

	return false
}

func (e Error) Wrap(err error) error {
	return wrappedError{
		cause: err,
		msg:   string(e),
	}
}

type wrappedError struct {
	cause error
	msg   string
}

func (w wrappedError) Error() string {
	if w.cause != nil {
		return fmt.Sprintf("%s%s%v", w.msg, ErrSeparator, w.cause)
	}

	return w.msg
}

func (w wrappedError) Is(err error) bool {
	return Error(w.msg).Is(err)
}

func (w wrappedError) As(err any) bool {
	return Error(w.msg).As(err)
}

func (w wrappedError) Unwrap() error {
	return w.cause
}

func New(message string) error {
	return errors.New(message)
}

func Is(err error, target error) bool {
	return errors.Is(err, target)
}

func As(err error, target any) bool {
	return errors.As(err, target)
}
