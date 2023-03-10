package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

// Table field type enum
const (
	int32Type int = iota
	uint32Type
	stringType
	floatType
	boolType

	int32VecType
	uint32VecType
	floatVecType
	stringVecType
	boolVecType

	int32Vec2Type
	uint32Vec2Type
	floatVec2Type
	stringVec2Type
	boolVec2Type
)

// type to its string
var type2String = map[int]string{
	int32Type:  "int32",
	uint32Type: "uint32",
	stringType: "string",
	floatType:  "float64",
	boolType:   "bool",

	int32VecType:  "[]int32",
	uint32VecType: "[]uint32",
	floatVecType:  "[]float64",
	stringVecType: "[]string",
	boolVecType:   "[]bool",

	int32Vec2Type:  "[][]int32",
	uint32Vec2Type: "[][]uint32",
	floatVec2Type:  "[][]float64",
	stringVec2Type: "[][]string",
	boolVec2Type:   "[][]bool",
}

type field struct {
	// field name.
	name_ string

	cvsField string
	// field type
	type_ string
}

type generator struct {
	typeName    string
	fileName    string  // in filepath
	filePathOut string  // out filepath
	fields      []field // file fields
}

func NewGenerator(outPath string) *generator {
	return &generator{
		filePathOut: outPath,
	}
}

// Get name.
func getName(filePath string) (string, string, bool) {
	data := strings.Split(filePath, "/")
	if len(data) <= 1 {
		return "", "", false
	}

	lastName := data[len(data)-1]
	data = strings.Split(lastName, ".")
	if len(data) <= 1 {
		return "", "", false
	} else {
		fileName := data[0]
		for i := 0; i < len(fileName); i++ {
			// Give up number or -, _ prefix tag.
			if !((fileName[i] >= '0' && fileName[i] <= '9') || fileName[i] == '_' || fileName[i] == '-' || fileName[i] == '&') {
				return fileName[i:], humpNaming(fileName[i:]), true
			}
		}
		// return error
		return "", "", false
	}
}

// aaa_bbb_cc return AaaBbbCc
func humpNaming(field string) string {
	if len(field) <= 0 {
		return field
	}

	var buf bytes.Buffer

	nextUpper := true

	for i := 0; i < len(field); i++ {
		if field[i] != '_' {
			if nextUpper && field[i] >= 'a' && field[i] <= 'z' {
				buf.WriteByte(field[i] - ('a' - 'A'))
				nextUpper = false
			} else {
				buf.WriteByte(field[i])
			}

		} else {
			nextUpper = true
		}
	}

	return buf.String()
}

// try to upper the first letter.
func upper(field string) string {
	if len(field) <= 0 {
		return field
	}

	if field[0] >= 'a' && field[0] <= 'z' { // if it lower alpha
		var buf bytes.Buffer
		c := field[0]
		buf.WriteByte(c - ('a' - 'A'))
		buf.WriteString(field[1:])
		return buf.String()
	} else {
		return field
	}
}

func hump(field string) string {
	if len(field) <= 0 {
		return field
	}
	split := strings.Split(field, "_")
	var humpTemp []string
	for _, s := range split {
		humpTemp = append(humpTemp, upper(s))
	}
	return strings.Join(humpTemp, "")
}

// try to lower the first letter.
func lower(field string) string {
	if len(field) <= 0 {
		return field
	}

	if field[0] >= 'A' && field[0] <= 'Z' { // if is upper alpha
		var buf bytes.Buffer
		c := field[0]
		buf.WriteByte(c + ('a' - 'A'))
		buf.WriteString(field[1:])
		return buf.String()
	} else {
		return field
	}
}

// load in file
func (g *generator) Load(filePath string) error {
	Parser := NewParser()
	if err := Parser.Load(filePath, true); err != nil {
		return err
	}

	// get file name, remove extension.
	var b bool
	if g.fileName, g.typeName, b = getName(filePath); !b {
		return errors.New("get file name error")
	}

	if Parser.GetAllCount() < 2 {
		return errors.New("csv format is error, it has now no remark info")
	}

	// read field name and its type
	fieldNameSize, _ := Parser.GetCount(0)
	g.fields = make([]field, 0, fieldNameSize)
	for i := 0; i < fieldNameSize; i++ {
		cvsField, _ := Parser.GetField(0, uint32(i))
		fieldName := upper(cvsField)
		fieldName = hump(fieldName)
		filedType, _ := Parser.GetField(1, uint32(i))
		fi := field{
			name_:    fieldName,
			cvsField: cvsField,
			type_:    filedType,
		}
		g.fields = append(g.fields, fi)
	}

	return nil
}

// Get name type: return string
func getNameType(Type string) string {
	logicType := getFieldType(Type)
	strType, ok := type2String[logicType]
	if ok {
		return strType
	} else {
		return "string"
	}
}

func (g *generator) buildInclude(buf *bytes.Buffer) {
	var con string
	con += "// ===================================== //\n"
	con += "// author:       gavingqf                //\n"
	con += "// == Please don'g change me by hand ==  //\n"
	con += "//====================================== //\n"
	con += "\n"

	con +=
		`/*you have defined the following interface:
type IConfig interface {
	// load interface
	Load(path string) bool

	// clear interface
	Clear()
}
*/`
	con += "\n\n"

	// add all imports as possible
	con += "package base\n\n"
	con += "import (\n"
	// con += "    \"log\"\n"
	con += "    \"shared/utility/glog\"\n"
	if g.checkArrayType() { // as golang is a strict language.
		con += "    \"strings\"\n\n"
	}
	// con += "    \"git.bilibili.co/gserver/log\"\n"
	con += ")\n\n"
	// more here

	buf.WriteString(con)
}

// Check whether exists array type(as strings.)
func (g *generator) checkArrayType() bool {
	for i := 0; i < len(g.fields); i++ {
		t := getFieldType(g.fields[i].type_)
		if t == floatVecType || t == int32VecType ||
			t == stringVecType || t == uint32VecType ||
			t == boolVecType || t == uint32Vec2Type ||
			t == int32Vec2Type || t == stringVec2Type ||
			t == floatVec2Type || t == boolVec2Type {
			return true
		}
	}
	return false
}

// just field alignment
func addSpace(field string, maxSize int) string {
	if len(field) >= maxSize {
		maxSize = len(field) + 2
	}

	var ret string
	for i := 0; i < maxSize-len(field); i++ {
		ret += " "
	}

	return ret
}

func (g *generator) buildStruct(buf *bytes.Buffer) {
	con := "type " + g.typeName + " struct {\n"
	maxFieldLen := 0
	for i := 0; i < len(g.fields); i++ {
		fieldLen := len(g.fields[i].name_)
		if maxFieldLen < fieldLen {
			maxFieldLen = fieldLen
		}
	}

	for i := 0; i < len(g.fields); i++ {
		con += addSpace("", 4) + g.fields[i].name_
		con += addSpace(g.fields[i].name_, maxFieldLen+1)
		con += getNameType(g.fields[i].type_)
		con += "\n"
	}
	con += "}\n\n"

	// write to file.
	buf.WriteString(con)
}

// keyWord check(id),here must be a id field.
func keyWordCheck(fieldName string) bool {
	return strings.EqualFold(fieldName, "id")
}

func (g *generator) getKeyType() string {
	for i := 0; i < len(g.fields); i++ {
		if keyWordCheck(g.fields[i].name_) {
			return getNameType(g.fields[i].type_)
		}
	}
	return "int32"
}

func (g *generator) buildConfig(buf *bytes.Buffer) bool {
	var con string

	con = "type " + g.typeName + "Config" + " struct {\n"
	con += addSpace("", 4)
	keyType := g.getKeyType()
	con += "data map[" + keyType + "]*" + g.typeName
	con += "\n"
	con += "}\n\n"

	// == New func ==
	con += "func New" + g.typeName + "Config() *" + g.typeName + "Config {\n"
	con += addSpace("", 4)
	con += "return &" + g.typeName + "Config{\n"
	con += addSpace("", 8)
	con += "data: make(map[" + keyType + "]*" + g.typeName + "),\n"
	con += addSpace("", 4)
	con += "}\n"
	con += "}\n\n"

	// == Load(...) func
	con += "func (c *" + g.typeName + "Config" + ") Load(filePath string) bool {\n"
	// read start
	con += addSpace("", 4)
	con += "parse := NewParser()\n"

	con += addSpace("", 4)
	con += "if err := parse.Load(filePath, true); err != nil {\n"
	con += addSpace("", 8)
	con += "glog.Info(\"Load\"" + ", filePath " + ",\"err: \", err)\n"
	con += addSpace("", 8)
	con += "return false\n"
	con += addSpace("", 4)
	con += "}\n\n"

	con += addSpace("", 4)
	con += "// iterator all lines' content"
	con += "\n"
	con += addSpace("", 4)
	con += "for i := 2; i < parse.GetAllCount(); i++ {\n"
	con += addSpace("", 8)
	con += "data := new(" + g.typeName + ")\n\n"

	// read field.
	var keyId string
	for i := 0; i < len(g.fields); i++ {
		field := "data." + g.fields[i].name_
		// maxSize := 20

		con += addSpace("", 8)
		switch getFieldType(g.fields[i].type_) {
		case stringType: // string
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			con += field + ", _ = parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				keyId = field // string can use as keyword.
			}

		case stringVecType: // []string
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \",\")"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "v := " + array + "[j]"
			con += "\n"
			con += addSpace("", 12)
			con += field + "= append(" + field + ", v)"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case stringVec2Type: // [][]string
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \";\")"
			con += "\n"
			con += addSpace("", 8)
			con += field + "= make(" + type2String[stringVec2Type] + ", len(" + array + "))"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "strValue := " + array + "[j]"
			con += "\n"

			con += addSpace("", 12)
			array2Variable := "vec2Value"
			con += array2Variable + " := strings.Split(strValue, \",\")"
			con += "\n"
			con += addSpace("", 12)
			con += "for n := 0; n < len(" + array2Variable + "); n++ {\n"
			con += addSpace("", 16)
			con += field + "[j] = append(" + field + "[j]," + array2Variable + "[n])"
			con += "\n"
			con += addSpace("", 12)
			con += "}\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case uint32Type: // uint32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			con += "v" + g.fields[i].name_ + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)
			con += "var " + g.fields[i].name_ + "Ret bool"
			con += "\n"
			con += addSpace("", 8)
			con += field + ", " + g.fields[i].name_ + "Ret" + " = String2Uint32(v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)

			// ret check
			con += "if !" + field + "Ret {"
			con += "\n"
			con += addSpace("", 12)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\",v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 12)
			con += "return false"
			con += addSpace("", 8)
			con += "}"

			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				keyId = field
			}

		case uint32VecType: // []uint32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			arrayVariable := "vec" + g.fields[i].name_
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += "if " + arrayVariable + " != \"\" {"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \",\")"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "v, ret := String2Uint32(" + array + "[j])"
			con += "\n"

			// ret check
			con += addSpace("", 12)
			con += "if !ret {"
			con += "\n"
			con += addSpace("", 16)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\"," + array + "[j])"
			con += "\n"
			con += addSpace("", 16)
			con += "return false"
			con += "\n"
			con += addSpace("", 12)
			con += "}"
			con += "\n"
			// add
			con += addSpace("", 12)
			con += field + "= append(" + field + ",v)"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case uint32Vec2Type: // [][]int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \";\")"
			con += "\n"
			con += addSpace("", 8)
			con += field + "= make(" + type2String[uint32Vec2Type] + ",len(" + array + "))"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "strValue := " + array + "[j]"
			con += "\n"

			con += addSpace("", 12)
			array2Variable := "vec2Value"
			con += array2Variable + " := strings.Split(strValue, \",\")"
			con += "\n"
			con += addSpace("", 12)
			con += "for n := 0; n < len(" + array2Variable + "); n++ {\n"
			con += addSpace("", 16)
			con += "v, _ := String2Uint32(" + array2Variable + "[n])\n"
			con += addSpace("", 16)
			con += field + "[j] = append(" + field + "[j], v)"
			con += "\n"
			con += addSpace("", 12)
			con += "}\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case int32Type: // int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			con += "v" + g.fields[i].name_ + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)
			con += "var " + g.fields[i].name_ + "Ret bool"
			con += "\n"
			con += addSpace("", 8)
			con += field + ", " + g.fields[i].name_ + "Ret = String2Int32(v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)

			// ret check
			con += "if !" + g.fields[i].name_ + "Ret {"
			con += "\n"
			con += addSpace("", 12)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\",v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 12)
			con += "return false"
			con += "\n"
			con += addSpace("", 8)
			con += "}"

			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				keyId = field
			}

		case int32VecType: // []int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += "if " + arrayVariable + " != \"\" {"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \",\")"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "v, ret := String2Int32(" + array + "[j])"
			con += "\n"

			// ret check
			con += addSpace("", 12)
			con += "if !ret {"
			con += "\n"
			con += addSpace("", 16)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error, value:\"," + array + "[j])"
			con += "\n"
			con += addSpace("", 16)
			con += "return false"
			con += "\n"
			con += addSpace("", 12)
			con += "}"
			con += "\n"
			// add
			con += addSpace("", 12)
			con += field + "= append(" + field + ", v)"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case int32Vec2Type: // [][]int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \";\")"
			con += "\n"
			con += addSpace("", 8)
			con += field + "= make(" + type2String[int32Vec2Type] + ",len(" + array + "))"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "strValue := " + array + "[j]"
			con += "\n"

			con += addSpace("", 12)
			array2Varible := "vec2Value"
			con += array2Varible + " := strings.Split(strValue, \",\")"
			con += "\n"
			con += addSpace("", 12)
			con += "for n := 0; n < len(" + array2Varible + "); n++ {\n"
			con += addSpace("", 16)
			con += "v, _ := String2Int32(" + array2Varible + "[n])\n"
			con += addSpace("", 16)
			con += field + "[j] = append(" + field + "[j], v)"
			con += "\n"
			con += addSpace("", 12)
			con += "}\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case floatType: // float32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			con += "v" + g.fields[i].name_ + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)
			con += "var " + g.fields[i].name_ + "Ret bool"
			con += "\n"
			con += addSpace("", 8)
			con += field + ", " + g.fields[i].name_ + "Ret = String2Float(v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)
			con += "if !" + g.fields[i].name_ + "Ret {"
			con += "\n"
			con += addSpace("", 12)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\",v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				keyId = field
			}

		case floatVecType: // []float32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \",\")"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "v, ret := String2Float(" + array + "[j])"
			con += "\n"

			// ret check
			con += addSpace("", 12)
			con += "if !ret {"
			con += "\n"
			con += addSpace("", 16)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\"," + array + "[j])"
			con += "\n"
			con += addSpace("", 16)
			con += "return false"
			con += "\n"
			con += addSpace("", 12)
			con += "}"
			con += "\n"
			// add
			con += addSpace("", 12)
			con += field + "= append(" + field + ",v)"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case floatVec2Type: // [][]int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \";\")"
			con += "\n"
			con += addSpace("", 8)
			con += field + " = make(" + type2String[floatVec2Type] + ", len(" + array + "))"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "strValue := " + array + "[j]"
			con += "\n"

			con += addSpace("", 12)
			array2Varible := "vec2Value"
			con += array2Varible + " := strings.Split(strValue, \",\")"
			con += "\n"
			con += addSpace("", 12)
			con += "for n := 0; n < len(" + array2Varible + "); n++ {\n"
			con += addSpace("", 16)
			con += "v, _ := String2Float(" + array2Varible + "[n])\n"
			con += addSpace("", 16)
			con += field + "[j] = append(" + field + "[j], v)"
			con += "\n"
			con += addSpace("", 12)
			con += "}\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case boolType: // bool
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			con += addSpace("", 8)
			con += "v" + g.fields[i].name_ + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)
			con += "var " + g.fields[i].name_ + "Ret bool"
			con += "\n"
			con += addSpace("", 8)
			con += field + ", " + g.fields[i].name_ + "Ret = String2Bool(v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)
			con += "if !" + g.fields[i].name_ + "Ret {"
			con += "\n"
			con += addSpace("", 12)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error,value:\",v" + g.fields[i].name_ + ")"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				keyId = field
			}

		case boolVecType: // []float32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \",\")"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "v, ret := String2Bool(" + array + "[j])"
			con += "\n"

			// ret check
			con += addSpace("", 12)
			con += "if !ret {"
			con += "\n"
			con += addSpace("", 16)
			con += "glog.Error(\"Parse " + g.typeName + "." + g.fields[i].name_ + " field error\")"
			con += "\n"
			con += addSpace("", 16)
			con += "return false"
			con += "\n"
			con += addSpace("", 12)
			con += "}"
			con += "\n"
			// add
			con += addSpace("", 12)
			con += field + "= append(" + field + ",v)"
			con += "\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}

		case boolVec2Type: // [][]int32
			con += "/* parse " + g.fields[i].name_ + " field */ \n"
			arrayVariable := "vec" + g.fields[i].name_
			con += addSpace("", 8)
			con += arrayVariable + ", _ := parse.GetFieldByName(uint32(i)," + "\"" + lower(g.fields[i].cvsField) + "\")"
			con += "\n"
			con += addSpace("", 8)

			array := "array" + upper(g.fields[i].name_)
			con += array + " := strings.Split(" + arrayVariable + ", \";\")"
			con += "\n"
			con += addSpace("", 8)
			con += field + "= make(" + type2String[boolVec2Type] + ", len(" + array + "))"
			con += "\n"
			con += addSpace("", 8)
			con += "for j := 0; j < len(" + array + "); j++ {\n"
			con += addSpace("", 12)
			con += "strValue := " + array + "[j]"
			con += "\n"

			con += addSpace("", 12)
			array2Varible := "vec2Value"
			con += array2Varible + " := strings.Split(strValue, \",\")"
			con += "\n"
			con += addSpace("", 12)
			con += "for n := 0; n < len(" + array2Varible + "); n++ {\n"
			con += addSpace("", 16)
			con += "v, _ := String2Bool(" + array2Varible + "[n])\n"
			con += addSpace("", 16)
			con += field + "[j] = append(" + field + "[j],v)"
			con += "\n"
			con += addSpace("", 12)
			con += "}\n"
			con += addSpace("", 8)
			con += "}"
			con += "\n\n"
			if keyWordCheck(g.fields[i].name_) {
				fmt.Println(g.typeName, " field is keyword, must be integer!!= ")
				return false
			}
		}
	}

	// repeated id check
	con += addSpace("", 8)
	con += "if _, ok := c.data[" + keyId + "]; ok {\n"
	con += addSpace("", 12)
	con += "glog.Errorf(\"Find %d repeated\"," + keyId + ")\n"
	con += addSpace("", 12)
	con += "return false"
	con += "\n"
	con += addSpace("", 8)
	con += "}\n"

	// add to map.
	con += addSpace("", 8)
	con += "c.data[" + keyId + "] = data\n"
	con += addSpace("", 4) + "}\n"

	// read end
	con += "\n"
	con += addSpace("", 4)
	con += "return true\n"
	con += "}"

	// == Clear() func
	con += "\n\n"
	con += "func (c *" + g.typeName + "Config" + ") Clear() {\n"
	con += "}\n\n"

	// == Find func
	con += "func (c *" + g.typeName + "Config" + ") Find(id " + g.getKeyType() + ") (*" + g.typeName + ", bool) {\n"
	con += addSpace("", 4)
	con += "v, ok := c.data[id]\n"
	con += addSpace("", 4)
	con += "return v, ok\n"
	con += "}\n\n"

	// == GetAllData func
	con += "func (c *" + g.typeName + "Config" + ") GetAllData() map[" + keyType + "]*" + g.typeName + "{\n"
	con += addSpace("", 4)
	con += "return c.data\n"
	con += "}\n\n"

	// == Traverse all data
	con += "func (c *" + g.typeName + "Config" + ") Traverse() {\n"
	con += addSpace("", 4)
	con += "for _, v := range c.data {\n"
	con += addSpace("", 8)
	str := ""
	for k := 0; k < len(g.fields); k++ {
		str += "v." + g.fields[k].name_
		if k == 0 {
			str += " "
		}
		if k != len(g.fields)-1 {
			str += ", \",\", "
		}
	}
	con += "glog.Info("
	con += str + ")\n"
	con += addSpace("", 4)
	con += "}\n"
	con += "}\n"
	buf.WriteString(con)

	return true
}

// Return filed type as enum
func getFieldType(typ string) int {
	typ = strings.Trim(typ, " ")

	switch {
	case strings.EqualFold(typ, "int"):
		return int32Type
	case strings.EqualFold(typ, "ints"):
		return int32VecType
	case strings.EqualFold(typ, "uint"):
		return uint32Type
	case strings.EqualFold(typ, "uints"):
		return uint32VecType
	case strings.EqualFold(typ, "string"):
		return stringType
	case strings.EqualFold(typ, "strings"):
		return stringVecType
	case strings.EqualFold(typ, "float"):
		return floatType
	case strings.EqualFold(typ, "floats"):
		return floatVecType
	case strings.EqualFold(typ, "bool"):
		return boolType
	case strings.EqualFold(typ, "bools"):
		return boolVecType
	case strings.EqualFold(typ, "int2"):
		return int32Vec2Type
	case strings.EqualFold(typ, "uint2"):
		return uint32Vec2Type
	case strings.EqualFold(typ, "string2"):
		return stringVec2Type
	case strings.EqualFold(typ, "float2"):
		return floatVec2Type
	case strings.EqualFold(typ, "bool2"):
		return boolVec2Type
	default:
		return stringType
	}
}

// building
func (g *generator) Build() bool {
	fileName := g.filePathOut + "/" + lower(g.fileName) + ".go"
	file, _ := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	defer file.Close()

	var buf bytes.Buffer

	g.buildInclude(&buf)
	g.buildStruct(&buf)
	if !g.buildConfig(&buf) {
		return false
	}

	// write all content to file
	_, err := file.Write(buf.Bytes())
	if err != nil {
		log.Printf("file %s write error: %v", fileName, err)
	}

	return true
}
