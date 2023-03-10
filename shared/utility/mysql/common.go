package mysql

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

func GetMysqlValue(value reflect.Value) ([]byte, error) {
	return getValue(value)
}

func SetMysqlValue(value reflect.Value, bs []byte) error {
	return setValue(value, bs)
}

func equal(v1, v2 []byte) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i := 0; i < len(v1); i++ {
		if v1[i] != v2[i] {
			return false
		}
	}

	return true
}

func getValue(value reflect.Value) ([]byte, error) {
	data := value.Interface()

	marshaler, ok := data.(Marshaler)
	if ok {
		content, err := marshaler.Marshal()
		if err != nil {
			return nil, err
		}

		return content, nil
	} else {
		switch value.Kind() {
		case reflect.Ptr, reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			return json.Marshal(data)
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int,
			reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
			reflect.Float32, reflect.Float64,
			reflect.Bool:

			return []byte(fmt.Sprint(data)), nil
		case reflect.String:
			return []byte(data.(string)), nil
		case reflect.Interface:
			_, ok := data.(string)
			if ok {
				return []byte(data.(string)), nil
			} else {
				return json.Marshal(data)
			}
		default:
			return []byte{}, nil
		}
	}
}

func setValue(value reflect.Value, bs []byte) error {
	if !value.CanSet() {
		return errors.New("value can't set")
	}

	data := value.Interface()

	marshaler, ok := data.(Marshaler)
	if ok {
		err := marshaler.Unmarshal(bs)
		if err != nil {
			return err
		}

		value.Set(reflect.ValueOf(data))
	} else {
		switch value.Kind() {
		case reflect.Ptr:
			err := json.Unmarshal(bs, data)
			if err != nil {
				return err
			}

			value.Set(reflect.ValueOf(data))
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			data = value.Addr().Interface()

			err := json.Unmarshal(bs, data)
			if err != nil {
				return err
			}

			value.Set(reflect.ValueOf(data).Elem())
		case reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int:
			x, err := strconv.ParseInt(string(bs), 10, 64)
			if err != nil {
				return nil
			}

			value.SetInt(x)
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
			x, err := strconv.ParseUint(string(bs), 10, 64)
			if err != nil {
				return nil
			}

			value.SetUint(x)
		case reflect.Float32, reflect.Float64:
			x, err := strconv.ParseFloat(string(bs), 64)
			if err != nil {
				return nil
			}

			value.SetFloat(x)
		case reflect.Bool:
			x, err := strconv.ParseBool(string(bs))
			if err != nil {
				return nil
			}

			value.SetBool(x)
		case reflect.String:
			value.SetString(string(bs))
		case reflect.Interface:
			_, ok := data.(string)
			if ok {
				value.SetString(string(bs))
			} else {
				newData := reflect.New(value.Type()).Interface()

				err := json.Unmarshal(bs, newData)
				if err != nil {
					return err
				}

				value.Set(reflect.ValueOf(newData).Elem())
			}
		default:
			return errors.New("invalid Value type")
		}
	}

	return nil
}

func GenObjectColumns(obj interface{}) string {
	buffer := &bytes.Buffer{}
	typ := reflect.TypeOf(obj).Elem()

	genReflectColumns(typ, buffer)

	return buffer.String()
}

func genReflectColumns(refType reflect.Type, buf *bytes.Buffer) int32 {
	if refType.Kind() != reflect.Struct {
		return 0
	}

	var colCount int32 = 0
	for i := 0; i < refType.NumField(); i++ {
		typeField := refType.Field(i)
		columnName := typeField.Tag.Get("json")
		ignore := typeField.Tag.Get("ignore") == "true"
		if columnName == "" && !ignore {
			colCount = colCount + genReflectColumns(typeField.Type.Elem(), buf)
			continue
		}

		if colCount > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('`')
		buf.WriteString(columnName)
		buf.WriteByte('`')
		colCount = colCount + 1
	}

	return colCount
}

func GenObjectValue(obj interface{}) (string, error) {
	buffer := &bytes.Buffer{}
	objVal := reflect.ValueOf(obj).Elem()
	_, err := genReflectValue(objVal, buffer)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func genReflectValue(refVal reflect.Value, buf *bytes.Buffer) (int32, error) {
	refType := refVal.Type()
	if refType.Kind() != reflect.Struct {
		return 0, nil
	}

	var colCount int32 = 0
	for i := 0; i < refType.NumField(); i++ {
		typeField := refType.Field(i)
		fieldVal := refVal.Field(i)

		fieldType := typeField.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldVal = fieldVal.Elem()
		}

		jsonTag := typeField.Tag.Get("json")
		ignore := typeField.Tag.Get("ignore") == "true"
		if jsonTag == "" && !ignore {
			addCol, err := genReflectValue(fieldVal, buf)
			colCount = colCount + addCol
			if err != nil {
				return colCount, err
			}
			continue
		}

		if colCount > 0 {
			buf.WriteByte(',')
		}

		bytes, err := GetMysqlValue(fieldVal)
		if err != nil {
			return colCount, err
		}
		if fieldVal.Type().Kind() == reflect.Bool {
			buf.Write(bytes)
		} else {
			buf.WriteByte('\'')
			buf.Write(bytes)
			buf.WriteByte('\'')
		}

		colCount = colCount + 1

	}

	return colCount, nil
}

func SetObjectValue(obj interface{}, values []interface{}) error {
	objVal := reflect.ValueOf(obj).Elem()
	var cursor int32 = 0
	return setReflectValue(objVal, values, &cursor)
}

func setReflectValue(refVal reflect.Value, values []interface{}, cursor *int32) error {
	refType := refVal.Type()
	if refType.Kind() != reflect.Struct {
		return nil
	}

	for i := 0; i < refType.NumField(); i++ {
		typeField := refType.Field(i)
		fieldVal := refVal.Field(i)

		fieldType := typeField.Type
		if fieldType.Kind() == reflect.Ptr {
			if fieldVal.IsNil() {
				fieldVal.Set(reflect.New(fieldVal.Type().Elem()))
			}
			fieldVal = fieldVal.Elem()
		}

		jsonTag := typeField.Tag.Get("json")
		ignore := typeField.Tag.Get("ignore") == "true"
		if jsonTag == "" && !ignore {
			err := setReflectValue(fieldVal, values, cursor)
			if err != nil {
				return err
			}
			continue
		}

		byteValue, ok := values[*cursor].(*[]byte)
		if !ok {
			return errors.New("invalid values type")
		}

		SetMysqlValue(fieldVal, *byteValue)

		*cursor = *cursor + 1
	}

	return nil
}

//////////////////////////
func GenInsertParams(obj interface{}) (columns, placeholders string, values []interface{}, err error) {
	columnBuffer := &bytes.Buffer{}
	placeholderBuffer := &bytes.Buffer{}
	refVal := reflect.ValueOf(obj).Elem()

	_, err = genInsertReflectParams(refVal, columnBuffer, placeholderBuffer, &values)
	if err != nil {
		return
	}

	columns = columnBuffer.String()
	placeholders = placeholderBuffer.String()

	return
}

func genInsertReflectParams(refVal reflect.Value, columnBuf, placeholderBuf *bytes.Buffer, values *[]interface{}) (int32, error) {
	refType := refVal.Type()
	if refType.Kind() != reflect.Struct {
		return 0, nil
	}

	var colCount int32 = 0
	for i := 0; i < refType.NumField(); i++ {
		typeField := refType.Field(i)
		fieldVal := refVal.Field(i)

		columnName := typeField.Tag.Get("db")
		if columnName == "-" {
			continue
		}
		if columnName == "" {
			addCol, err := genInsertReflectParams(fieldVal, columnBuf, placeholderBuf, values)
			colCount = colCount + addCol
			if err != nil {
				return colCount, err
			}
			continue
		}

		if colCount > 0 {
			columnBuf.WriteByte(',')
			placeholderBuf.WriteByte(',')
		}
		columnBuf.WriteByte('`')
		columnBuf.WriteString(columnName)
		columnBuf.WriteByte('`')

		placeholderBuf.WriteString("?")

		valBytes, err := GetMysqlValue(fieldVal)
		if err != nil {
			return colCount, err
		}

		*values = append((*values), string(valBytes))

		colCount = colCount + 1
	}

	return colCount, nil
}

func GenUpdateParams(obj interface{}) (params string, values []interface{}, err error) {
	paramsBuffer := &bytes.Buffer{}
	refVal := reflect.ValueOf(obj).Elem()

	_, err = genUpdateReflectParams(refVal, paramsBuffer, &values)
	if err != nil {
		return
	}

	params = paramsBuffer.String()

	return
}

func genUpdateReflectParams(refVal reflect.Value, paramBuf *bytes.Buffer, values *[]interface{}) (int32, error) {
	refType := refVal.Type()
	if refType.Kind() != reflect.Struct {
		return 0, nil
	}

	var colCount int32 = 0
	for i := 0; i < refType.NumField(); i++ {
		typeField := refType.Field(i)
		fieldVal := refVal.Field(i)

		columnName := typeField.Tag.Get("db")
		if columnName == "-" {
			continue
		}
		if columnName == "" {
			addCol, err := genUpdateReflectParams(fieldVal, paramBuf, values)
			colCount = colCount + addCol
			if err != nil {
				return colCount, err
			}
			continue
		}

		if colCount > 0 {
			paramBuf.WriteByte(',')
		}
		paramBuf.WriteByte('`')
		paramBuf.WriteString(columnName)
		paramBuf.WriteByte('`')
		paramBuf.WriteString("=?")

		valBytes, err := GetMysqlValue(fieldVal)
		if err != nil {
			return colCount, err
		}

		*values = append((*values), string(valBytes))

		colCount = colCount + 1
	}

	return colCount, nil
}
