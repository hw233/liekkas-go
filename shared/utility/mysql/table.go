package mysql

import (
	"reflect"

	"shared/utility/naming"
)

// TableStruct的FieldStruct和FieldValue是必须一一对应的，改变前者长度会同时影响后者长度，后者不允许改变长度
type Table struct {
	*TableStruct
	values []*FieldValue
}

func NewTable(ts *TableStruct) *Table {
	table := &Table{
		TableStruct: ts,
		values:      make([]*FieldValue, 0, ts.FieldsNum()),
	}

	for i := 0; i < ts.FieldsNum(); i++ {
		table.values = append(table.values, NewInitFieldValue())
	}

	return table
}

// -------------------------------- FieldValue ---------------------------------

func (t *Table) Field(i int) (*Field, error) {
	if i >= t.FieldsNum() {
		return nil, ErrFieldIndexOutOfRange
	}

	// 上面判断过了，这个不需要判断了
	fieldStruct, _ := t.TableStruct.FieldStruct(i)

	return NewField(fieldStruct, t.values[i]), nil
}

func (t *Table) MajorField() *Field {
	field, _ := t.Field(t.MajorIndex)
	return field
}

func (t *Table) FieldValue(i int) (*FieldValue, error) {
	if i >= t.FieldsNum() {
		return nil, ErrFieldIndexOutOfRange
	}

	return t.values[i], nil
}

func (t *Table) SetFields(fields []*FieldStruct) {
	t.TableStruct.SetFields(fields)

	l := len(fields)

	values := make([]*FieldValue, 0, l)

	for i := 0; i < l; i++ {
		values = append(values, NewInitFieldValue())
	}

	t.values = values
}

func (t *Table) AddFields(name string, index int, typ reflect.Type) {
	t.TableStruct.AddFields(name, index, typ)
	t.values = append(t.values, NewInitFieldValue())
}

func (t *Table) AppendFields(field *FieldStruct) {
	t.TableStruct.AppendFields(field)
	t.values = append(t.values, NewInitFieldValue())
}

func (t *Table) NeedUpdate() bool {
	for i, value := range t.values {
		if t.MajorIndex != i && value.NeedUpdate {
			return true
		}
	}

	return false
}

func (t *Table) ReflectValue(data interface{}) error {
	ve := reflect.ValueOf(data).Elem()

	for i := 0; i < t.FieldsNum(); i++ {
		field, err := t.Field(i)
		if err != nil {
			return err
		}

		// update Field
		field.Update(ve.Field(field.Index).Interface())
	}

	return nil
}

func (t *Table) ReflectSet(data interface{}) error {
	ve := reflect.ValueOf(data).Elem()

	for i := 0; i < t.FieldsNum(); i++ {
		field, err := t.Field(i)
		if err != nil {
			return err
		}

		err = setValue(ve.Field(field.Index), field.Value)
		if err != nil {
			return err
		}
	}

	return nil
}

type TableStruct struct {
	typ        reflect.Type
	Name       string
	MajorIndex int
	Fields     []*FieldStruct
}

func NewTableStruct(name string) *TableStruct {
	return &TableStruct{
		typ:        nil,
		Name:       name,
		MajorIndex: 0,
		Fields:     []*FieldStruct{},
	}
}

func (t *TableStruct) IsNew() bool {
	return t.Name == "" || t.Fields == nil
}

func (t *TableStruct) FieldStruct(i int) (*FieldStruct, error) {
	if i >= t.FieldsNum() {
		return nil, ErrFieldIndexOutOfRange
	}

	return t.Fields[i], nil
}

func (t *TableStruct) SetFields(fields []*FieldStruct) {
	t.Fields = fields
}

func (t *TableStruct) AddFields(name string, index int, typ reflect.Type) {
	t.Fields = append(t.Fields, NewFieldStruct(name, index, typ))
}

func (t *TableStruct) AppendFields(field *FieldStruct) {
	t.Fields = append(t.Fields, field)
}

func (t *TableStruct) FieldsNum() int {
	return len(t.Fields)
}

func (t *TableStruct) MajorFieldName() string {
	return t.Fields[t.MajorIndex].Name
}

func (t *TableStruct) AssignableTo(data interface{}) bool {
	if t.typ == nil {
		return false
	}
	return t.typ.AssignableTo(reflect.TypeOf(data))
}

func (t *TableStruct) ReflectTable(data interface{}) error {
	t.typ = reflect.TypeOf(data)
	te := t.typ.Elem()

	majorIndex := -1

	for i := 0; i < te.NumField(); i++ {
		sft := te.Field(i)
		tag := NewTag(sft.Tag)

		if tag.DB == "-" {
			continue
		}

		name := ""

		if tag.DB != "" {
			name = tag.DB
		} else {
			name = naming.UnderlineNaming(sft.Name)
		}

		if tag.Major == "true" {
			if majorIndex != -1 {
				return ErrOvermuchMajor
			}

			majorIndex = i
		}

		// set Field
		t.AddFields(name, i, sft.Type)
	}

	// no db field
	if t.FieldsNum() < 1 {
		return ErrNoField
	}

	// no major key
	if majorIndex == -1 {
		return ErrNotFoundMajor
	}

	t.MajorIndex = majorIndex

	return nil
}
