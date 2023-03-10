package common

import (
	"reflect"
)

type GMTask struct {
	f  reflect.Value
	in []reflect.Value
}

func NewGMTask(f reflect.Value, in ...interface{}) *GMTask {
	inVal := make([]reflect.Value, 0, len(in))
	for i, _ := range in {
		inVal = append(inVal, reflect.ValueOf(in[i]))
	}

	return &GMTask{
		f:  f,
		in: inVal,
	}
}

func (t *GMTask) Do() error {
	ret := t.f.Call(t.in)
	if len(ret) > 0 {
		err, ok := ret[0].Interface().(error)
		if ok {
			return err
		}
	}

	return nil
}
