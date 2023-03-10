package errors

import (
	"reflect"
	"runtime"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	defaultCode  = 0
	defaultTitle = "(/ﾟДﾟ)/ A Error Occurred!"
	textDivide   = "; "

	_code  = "Code"
	_text  = "Text"
	_trace = "Trace"
)

func SetDefaultCode(code int) {
	defaultCode = code
}

func SetTitle(title string) {
	defaultTitle = title
}

func SetTextDivide(d string) {
	textDivide = d
}

func cause(err error) *baseError {
	return &baseError{err}
}

func caller() uintptr {
	pc, _, _, _ := runtime.Caller(3)
	return pc
}

func New(format string, v ...interface{}) error {
	return newPCError(newTextError(nil, format, v...))
}

func NewCode(code int, text string) error {
	return newCodeError(newPCError(newTextError(nil, text)), code)
}

func Wrap(err error, text string, v ...interface{}) error {
	return newPCError(newTextError(newDivideError(err), text, v...))
}

func Swrapf(err error, v ...interface{}) error {
	if Code(err) != defaultCode {
		return newCodeError(newPCError(newTextError(nil, Text(err), v...)), Code(err))
	}

	return newPCError(newTextError(nil, Text(err), v...))
}

func WrapText(err error, text string, v ...interface{}) error {
	return newTextError(newDivideError(err), text, v...)
}

func WrapTrace(err error) error {
	return newPCError(newDivideError(err))
}

func localizeError(err error) error {
	if _, ok := err.(interface {
		WrapError(error)
	}); !ok && err != nil {
		return newTextError(nil, err.Error())
	}

	return err
}

func WrapError(err, werr error) error {
	if werr == nil {
		return err
	}

	err, werr = localizeError(err), localizeError(werr)

	w := werr

	for err != nil {
		if u, ok := w.(interface {
			Unwrap() error
		}); ok {
			if cause := u.Unwrap(); cause != nil {
				w = cause

				continue
			}
		}

		we := w.(interface {
			WrapError(error)
		})

		we.WrapError(newDivideError(err))

		break
	}

	return werr
}

func Unwrap(err error) error {
	for {
		if u, ok := err.(interface {
			Unwrap() error
		}); ok {
			if i, ok := err.(interface {
				IsDivide() bool
			}); ok {
				if i.IsDivide() {
					return u.Unwrap()
				}

				err = u.Unwrap()

				continue
			}
		}

		break
	}

	return nil
}

func Cause(err error) error {
	for {
		if cause := Unwrap(err); cause != nil {
			err = cause
			continue
		}

		break
	}

	return err
}

func Is(err, target error) bool {
	if err == nil || target == nil {
		return err == target
	}

	for {
		if reflect.DeepEqual(err, target) {
			return true
		}

		if x, ok := err.(interface {
			Is(error) bool
		}); ok && x.Is(target) {
			return true
		}

		if err = Unwrap(err); err == nil {
			return false
		}
	}
}

func Format(err error) *fmtError {
	return newFmtError(err)
}

func Code(err error) int {
	if code := code(err); code != defaultCode {
		return code
	}

	s, ok := status.FromError(err)
	if ok && s != nil && s.Code() != codes.OK {
		return int(s.Code())
	}

	return defaultCode
}

func code(err error) int {
	if e, ok := err.(interface {
		Code() int
	}); ok {
		return e.Code()
	}

	return defaultCode
}

func Text(err error) string {
	err = cause(err)

	// s, ok := status.FromError(err)
	// if ok {
	// 	return s.Message()
	// }

	if e, ok := err.(interface {
		Text() string
	}); ok {
		return e.Text()
	}

	return ""
}

func text(err error) string {
	if e, ok := err.(interface {
		Text() string
	}); ok {
		return e.Text()
	}

	return ""
}

func Trace(err error) string {
	err = cause(err)

	s, ok := status.FromError(err)
	if ok && s != nil {
		return s.Message()
	}

	return trace(err)
}

func trace(err error) string {
	if e, ok := err.(interface {
		Trace() string
	}); ok {
		return e.Trace()
	}

	return ""
}
