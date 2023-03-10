package errors

import (
	"fmt"
)

type textError struct {
	*baseError
	text string
}

func newTextError(err error, text string, v ...interface{}) *textError {
	t := ""

	if len(v) != 0 {
		t = fmt.Sprintf(text, v...)
	} else {
		t = text
	}

	return &textError{
		baseError: cause(err),
		text:      t,
	}
}

func (e *textError) Text() string {
	return e.text
}

func (e *textError) Error() string {
	if e.cause != nil {
		if ce := e.cause.Error(); ce != "" {
			return ce + textDivide + e.text
		}
	}

	return e.text
}

func (e *textError) Trace() string {
	return e.text + "\n"
}

// func (e *textError) Swrapf(v ...interface{}) {
// 	e.text = fmt.Sprintf(e.text, v...)
// }
