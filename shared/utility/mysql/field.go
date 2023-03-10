package mysql

import (
	"reflect"
)

type Field struct {
	*FieldStruct
	*FieldValue
}

func NewField(fieldStruct *FieldStruct, fieldValue *FieldValue) *Field {
	return &Field{
		FieldStruct: fieldStruct,
		FieldValue:  fieldValue,
	}
}

type FieldValue struct {
	Value      []byte
	NeedUpdate bool
}

func NewInitFieldValue() *FieldValue {
	return &FieldValue{
		Value:      []byte{},
		NeedUpdate: false,
	}
}

// func NewFieldValue(value []byte) *FieldValue {
// 	// var v []byte
// 	// var err error
//
// 	// rv, ok := value.(reflect.Value)
// 	// if ok {
// 	// 	v, err = getValue(rv)
// 	// 	if err != nil {
// 	// 		v = []byte{}
// 	// 	}
// 	// } else {
// 	// 	v, err = getValue(reflect.ValueOf(value))
// 	// 	if err != nil {
// 	// 		v = []byte{}
// 	// 	}
// 	// }
//
// 	return &FieldValue{
// 		Value:      value,
// 		NeedUpdate: false,
// 	}
// }

func (fv *FieldValue) ValuePtr() *[]byte {
	return &fv.Value
}

func (fv *FieldValue) Interface() interface{} {
	var i interface{}
	_ = setValue(reflect.ValueOf(i), fv.Value)
	return i
}

func (fv *FieldValue) Update(value interface{}) {
	var v []byte
	var err error

	rv, ok := value.(reflect.Value)
	if ok {
		v, err = getValue(rv)
		if err != nil {
			v = []byte{}
		}
	} else {
		v, err = getValue(reflect.ValueOf(value))
		if err != nil {
			v = []byte{}
		}
	}

	if !equal(fv.Value, v) {
		fv.NeedUpdate = true
		fv.Value = v
	}
}

func (fv *FieldValue) RefreshUpdate() {
	fv.NeedUpdate = false
}

type FieldStruct struct {
	Name  string
	Index int
	Type  reflect.Type
}

func NewFieldStruct(name string, index int, typ reflect.Type) *FieldStruct {
	return &FieldStruct{
		Name:  name,
		Index: index,
		Type:  typ,
	}
}

func (s *FieldStruct) NewValue() interface{} {
	return reflect.New(s.Type).Interface()
}
