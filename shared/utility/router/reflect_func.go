package router

import "reflect"

type ReflectFunc struct {
	f  reflect.Value
	in []reflect.Type
}

func NewReflectFunc(f interface{}) *ReflectFunc {
	t := reflect.TypeOf(f)

	if t.Kind() != reflect.Func {
		return nil
	}

	in := make([]reflect.Type, 0, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		in = append(in, t.In(i))
	}

	return &ReflectFunc{
		f:  reflect.ValueOf(f),
		in: in,
	}
}

func (rf *ReflectFunc) Call(in ...interface{}) []interface{} {
	inVal := make([]reflect.Value, 0, len(in))

	for i := range in {
		inVal = append(inVal, reflect.ValueOf(in[i]))
	}

	var out []interface{}

	outVal := rf.f.Call(inVal)

	for i := range outVal {
		out = append(out, outVal[i].Interface())
	}

	return out
}

func (rf *ReflectFunc) In(i int) interface{} {
	if i >= len(rf.in) {
		return nil
	}

	return reflect.New(rf.in[i].Elem()).Interface()
}

func (rf *ReflectFunc) NumIn() int {
	return len(rf.in)
}

func (rf *ReflectFunc) Interface() interface{} {
	return rf.f.Interface()
}
