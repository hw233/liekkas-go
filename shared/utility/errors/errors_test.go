package errors

import (
	"errors"
	"testing"
)

// Wrap, UnWarp and Cause
func TestCover(t *testing.T) {
	checkErrorMsg := func(err error, msg string) {
		if m := err.Error(); m != msg {
			t.Errorf("error text: %s, want: %s", m, msg)
		}
	}

	checkErrorEqual := func(err1, err2 error) {
		if err1 != err2 {
			t.Errorf("error: %v, want: %v", err1, err2)
		}
	}

	// nil error
	var nilErr error

	// go error
	goErr := errors.New("go error")

	// error without code
	msgErr := New("error")

	// error with code
	codeErr := NewCode(233, "code error")

	// Wrap
	checkErrorMsg(WrapText(nilErr, "warp text"), "warp text")
	checkErrorMsg(WrapText(goErr, "warp text"), "go error; warp text")
	checkErrorMsg(WrapText(msgErr, "warp text"), "error; warp text")

	wrapCodeErr := Wrap(codeErr, "warp message1")
	checkErrorMsg(wrapCodeErr, "code error; warp message1")
	wrapWrapCodeErr := Wrap(wrapCodeErr, "warp message2")
	checkErrorMsg(wrapWrapCodeErr, "code error; warp message1; warp message2")

	// Unwrap
	checkErrorEqual(Unwrap(nilErr), nil)
	checkErrorEqual(Unwrap(goErr), nil)
	checkErrorEqual(Unwrap(msgErr), nil)
	checkErrorEqual(Unwrap(wrapWrapCodeErr), wrapCodeErr)

	// Cause
	checkErrorEqual(Cause(nilErr), nil)
	checkErrorEqual(Cause(goErr), goErr)
	checkErrorEqual(Cause(msgErr), msgErr)
	checkErrorEqual(Cause(wrapWrapCodeErr), codeErr)
}

func TestIs(t *testing.T) {
	checkIs := func(err1, err2 error, ret bool) {
		if is := Is(err1, err2); is != ret {
			t.Errorf("err1: %v, err2: %v, is: %v, want: %v", err1, err2, is, ret)
		}
	}

	// nil error
	var nilErr error

	// go error
	goErr := errors.New("go error")
	wrapGoErr := Wrap(goErr, "warp text")

	// error without code
	msgErr := New("go error")

	// error with code
	codeErr1 := NewCode(233, "code error1")
	codeErr2 := NewCode(666, "code error1") // text equals to codeErr1
	codeErr3 := NewCode(233, "code error3") // code equals to codeErr1
	codeErr4 := NewCode(0, "go error")      // text equals to goErr

	checkIs(nilErr, nilErr, true)
	checkIs(goErr, errors.New("go error"), true) // same text

	checkIs(nilErr, goErr, false)
	checkIs(goErr, nilErr, false)
	checkIs(goErr, msgErr, false)      // same text but code different, goErr.code = nil
	checkIs(codeErr1, codeErr2, false) // same text but code different

	checkIs(wrapGoErr, goErr, true)
	checkIs(codeErr1, codeErr3, true) // same code
	checkIs(Wrap(codeErr1, "warp text"), codeErr1, true)
	checkIs(WrapError(codeErr1, msgErr), codeErr1, true)

	checkIs(goErr, codeErr4, false) // goErr not implements Is()
	checkIs(codeErr4, goErr, false) // codeErr4 implements Is() but equal to defaultCode
	checkIs(msgErr, goErr, false)   // msgErr implements Is() but msgErr.code = nil

	checkIs(WrapError(codeErr1, codeErr2), codeErr1, true) // code changes but UnWarp(err) equal to codeErr1
	checkIs(WrapError(codeErr3, codeErr2), codeErr1, true) // code changes but code of UnWarp(err) equal to codeErr1
}

func TestErrorCode(t *testing.T) {
	errorCode1, errorCode2, errorCode3 := 233, 666, 1024

	checkErrorCode := func(err error, code int) {
		if c := Code(err); c != code {
			t.Errorf("error code: %d, want: %d", c, code)
		}
	}

	SetDefaultCode(0)

	// New
	// nil error
	var nilErr error
	checkErrorCode(nilErr, 0)

	// go error
	goErr := errors.New("go error")
	checkErrorCode(goErr, 0)

	// error without code
	msgErr := New("error")
	checkErrorCode(msgErr, 0)

	// error with code
	codeErr := NewCode(errorCode2, "code error")
	checkErrorCode(codeErr, errorCode2)

	// set default error code and check again
	SetDefaultCode(errorCode1)
	checkErrorCode(nilErr, errorCode1)
	checkErrorCode(goErr, errorCode1)
	checkErrorCode(msgErr, errorCode1)
	checkErrorCode(NewCode(0, "code error"), 0)

	// Wrap
	checkErrorCode(Wrap(nilErr, "wrap text"), errorCode1)
	checkErrorCode(Wrap(goErr, "wrap text"), errorCode1)
	checkErrorCode(Wrap(msgErr, "wrap text"), errorCode1)
	checkErrorCode(Wrap(codeErr, "wrap text"), errorCode2)

	// WrapError
	checkErrorCode(WrapError(nilErr, nilErr), errorCode1)
	checkErrorCode(WrapError(goErr, goErr), errorCode1)
	checkErrorCode(WrapError(nilErr, codeErr), errorCode2)
	checkErrorCode(WrapError(codeErr, nilErr), errorCode2)
	checkErrorCode(WrapError(goErr, codeErr), errorCode2)
	checkErrorCode(WrapError(codeErr, goErr), errorCode2)
	checkErrorCode(WrapError(codeErr, NewCode(errorCode3, "other code error")), errorCode3)
}

func TestSwrapf(t *testing.T) {
	err := NewCode(233, "code errors %s:%d")
	err1 := Swrapf(err, "222222", 222222)
	err2 := Swrapf(err, "222222", 222222)
	t.Logf("%+v", err1)
	t.Logf("%+v", err2)
}

func TestFormat(t *testing.T) {
	err := NewCode(233, "code errors %s:%d")
	err = Swrapf(err, "111111", 111111)
	// err = WrapText(err, "wrap text %s:%d")
	// err = Swrapf(err, "222222", 222222)
	err = WrapTrace(err)
	err = WrapTrace(err)
	err = Wrap(err, "wrap text and trace")
	err = WrapError(err, errors.New("go error"))
	err = WrapError(err, NewCode(666, "code error2"))
	err = WrapError(err, New("errorf %s", "fffff"))
	err = Format(err).AddAttr("UID", 123456).
		AddAttr("OpenID", "xxxxxxxx").
		Rename("Code", "Code_New").
		SetTitle("New Title").
		SetOrder([]string{"UID", "OpenID", "Message", "Code_New", "Trace"})
	t.Logf("%+v", err)
}
