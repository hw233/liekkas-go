package errors

import "fmt"

type fmtError struct {
	*baseError
	title string
	fmts  map[string]interface{}
	order []string
}

func newFmtError(err error) *fmtError {
	fe := &fmtError{
		baseError: cause(err),
		title:     defaultTitle,
		fmts:      nil,
		order:     []string{_code, _text, _trace},
	}

	fe.fmts = map[string]interface{}{
		_code:  fe.Code,
		_text:  fe.Error,
		_trace: fe.Trace,
	}

	return fe
}

func (e *fmtError) Code() int {
	return Code(e.cause)
}

func (e *fmtError) Trace() string {
	return fmt.Sprintf("%-v", e)
}

func (e *fmtError) AddAttr(key string, val interface{}) *fmtError {
	e.fmts[key] = val
	if e.order[len(e.order)-1] == _trace {
		e.order = append(e.order[:len(e.order)-1], key, e.order[len(e.order)-1])
		return e
	}

	e.order = append(e.order, key)

	return e
}

func (e *fmtError) Rename(oldName, newName string) *fmtError {
	switch oldName {
	case _code:
		_code = newName
	case _text:
		_text = newName
	case _trace:
		_trace = newName
	}

	if v, ok := e.fmts[oldName]; ok {
		e.fmts[newName] = v
		delete(e.fmts, oldName)
	}

	for i, v := range e.order {
		if v == oldName {
			e.order[i] = newName
		}
	}

	return e
}

func (e *fmtError) SetOrder(order []string) *fmtError {
	e.order = order
	return e
}

func (e *fmtError) SetTitle(title string) *fmtError {
	e.title = title
	return e
}

func (e *fmtError) Format(s fmt.State, verb rune) {
	str := e.Error()

	switch verb {
	case 'v':
		switch {
		case s.Flag('+'):
			str = e.title + "\n"

			for _, o := range e.order {
				if v, ok := e.fmts[o]; ok {
					switch o {
					case _code:
						if f, ok := v.(func() int); ok {
							str += fmt.Sprintf("%s: %d\n", o, f())
						}
					case _text:
						if f, ok := v.(func() string); ok {
							str += fmt.Sprintf("%s: %s\n", o, f())
						}
					case _trace:
						if f, ok := v.(func() string); ok {
							str += fmt.Sprintf("%s:\n%s", o, f())
						}
					default:
						str += fmt.Sprintf("%s: %v\n", o, v)
					}
				}
			}
		case s.Flag('-'):
			// if ce, ok := e.cause.(interface {
			// 	Trace() string
			// }); ok {
			// 	str = ce.Trace()
			// }

			str = trace(e.baseError)

			if e.cause != nil {
				str = fmt.Sprintf("%-v", e.cause) + str
			}
		}
	}

	_, _ = s.Write([]byte(str))
}
