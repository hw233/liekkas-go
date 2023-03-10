package errors

type divideError struct {
	*baseError
}

func newDivideError(err error) *divideError {
	return &divideError{
		baseError: cause(err),
	}
}

func (e *divideError) IsDivide() bool {
	return true
}

func (e *divideError) Text() string {
	return ""
}

func (e *divideError) Trace() string {
	return ""
}
