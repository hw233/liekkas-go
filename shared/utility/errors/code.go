package errors

type codeError struct {
	*baseError
	code int
}

func newCodeError(err error, code int) *codeError {
	return &codeError{
		baseError: cause(err),
		code:      code,
	}
}

func (e *codeError) Code() int {
	return e.code
}

func (e *codeError) Trace() string {
	return ""
}

func (e *codeError) Is(err error) bool {
	return e.code != defaultCode && e.code == Code(err)
}
