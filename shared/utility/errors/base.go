package errors

import (
	"fmt"

	"google.golang.org/grpc/status"
)

type baseError struct {
	cause error
}

func (e *baseError) Format(s fmt.State, verb rune) {
	str := e.Error()

	switch verb {
	case 'v':
		switch {
		case s.Flag('-'):
			str = trace(e.cause)

			if e.cause != nil {
				str = fmt.Sprintf("%-v", e.cause) + str
			}
		}
	}

	_, _ = s.Write([]byte(str))
}

// func (e *baseError) Trace() string {
// 	return ""
// }

// func (e *baseError) trace() string {
// 	// s, ok := status.FromError(e.cause)
// 	// if ok {
// 	// 	return s.Message()
// 	// }
//
// 	if e, ok := e.cause.(interface {
// 		Trace() string
// 	}); ok {
// 		return e.Trace()
// 	}
//
// 	return ""
// }

func (e *baseError) Code() int {
	if e.cause == nil {
		return defaultCode
	}

	return Code(e.cause)
}

func (e *baseError) Text() string {
	if e.cause == nil {
		return ""
	}

	return text(e.cause)
}

// func (e *baseError) Swrapf(v ...interface{}) {
// 	if e, ok := e.cause.(interface {
// 		Swrapf(...interface{})
// 	}); ok {
// 		e.Swrapf(v...)
// 	}
// }

// func (e *baseError) Text() string {
// 	glog.Info("-------------->base _text")
// 	if e, ok := e.cause.(interface {
// 		Text() string
// 	}); ok {
// 		return e.Text()
// 	}
//
// 	return ""
// }

func (e *baseError) Is(err error) bool {
	if e, ok := e.cause.(interface {
		Is(error) bool
	}); ok {
		return e.Is(err)
	}

	return false
}

func (e *baseError) Error() string {
	if e.cause == nil {
		return ""
	}

	return e.cause.Error()
}

func (e *baseError) WrapError(err error) {
	if e.cause == nil {
		e.cause = err
	}
}

func (e *baseError) Unwrap() error {
	return e.cause
}

func (e *baseError) IsDivide() bool {
	return false
}

func (e *baseError) GRPCStatus() *status.Status {
	s, ok := status.FromError(e.cause)
	if s == nil || !ok {
		return nil
		// return status.New(codes.Code(code(e)), trace(e))
	}

	// if code := code(e); code != 0 {
	// 	return status.New(codes.Code(code), s.Message()+trace(e.cause))
	// }

	return status.New(s.Code(), s.Message()+trace(e.cause))
}
