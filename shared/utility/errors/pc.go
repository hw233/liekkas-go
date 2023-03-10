package errors

import (
	"fmt"
	"runtime"
)

type pcError struct {
	*baseError
	pc uintptr
}

func newPCError(err error) *pcError {
	return &pcError{
		baseError: cause(err),
		pc:        caller(),
	}
}

func (e *pcError) Trace() string {
	fc := runtime.FuncForPC(e.pc)
	file, line := fc.FileLine(e.pc)

	return fmt.Sprintf("\t%s()\n\t\t%s:%d\t\n", fc.Name(), file, line)
}
