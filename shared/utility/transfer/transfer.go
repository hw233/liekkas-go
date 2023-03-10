package transfer

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

const (
	// TODO: 后续根据需求可以添加支持的类型
	tagSrc    = "src"    // 数据源，多个数据源用逗号隔开，支持数组，切片，例：Name,Slice[1],Slice[2]
	tagRule   = "rule"   // 用指定转换规则去转换
	tagIgnore = "ignore" // 忽略该字段

	ruleDefault = "_default" // 默认规则
)

var (
	lock           sync.RWMutex
	rules          = map[string]interface{}{}
	isRulesChecked bool
)

func RegisterRule(name string, rule interface{}) {
	lock.Lock()
	defer lock.Unlock()

	rules[name] = rule
	isRulesChecked = false
}

// 每次注册检测错误处理起来会比较难看，所以注册完务必检测一遍，否则为了防止崩溃，使用时候会自动检查一次
func CheckRules() error {
	lock.Lock()
	defer lock.Unlock()

	for _, rule := range rules {
		ruleType := reflect.TypeOf(rule)
		if ruleType.Kind() != reflect.Func {
			return fmt.Errorf("rule %s not func", rule)
		}

		// out 1~2
		if ruleType.NumOut() < 1 {
			return fmt.Errorf("rule %s not out", rule)
		}

		if ruleType.NumOut() > 2 {
			return fmt.Errorf("rule %s too much out", rule)
		}
	}

	isRulesChecked = true

	return nil
}

func getRule(name string) (interface{}, bool) {
	lock.Lock()
	defer lock.Unlock()

	rule, ok := rules[name]
	return rule, ok
}

func Transfer(src, dst interface{}) error {
	if !isRulesChecked {
		err := CheckRules()
		if err != nil {
			return err
		}
	}

	srcTypePtr := reflect.TypeOf(src)
	var srcType reflect.Type
	if srcTypePtr.Kind() == reflect.Ptr {
		srcType = srcTypePtr.Elem()
		if srcType.Kind() != reflect.Struct {
			return errors.New("type of src isn't ptr of struct")
		}
	} else if srcTypePtr.Kind() == reflect.Struct {
		srcType = srcTypePtr
	} else {
		return errors.New("type of src isn't struct or ptr of struct")
	}

	srcValue := reflect.ValueOf(src).Elem()

	dstTypePtr := reflect.TypeOf(dst)
	if dstTypePtr.Kind() != reflect.Ptr {
		return errors.New("type of dst isn't ptr")
	}

	dstType := dstTypePtr.Elem()
	if dstType.Kind() != reflect.Struct {
		return errors.New("type of dst isn't ptr of struct")
	}

	dstValue := reflect.ValueOf(dst).Elem()

	values := map[string]reflect.Value{}

	for i := 0; i < srcType.NumField(); i++ {
		values[srcType.Field(i).Name] = srcValue.Field(i)
	}

	for i := 0; i < dstType.NumField(); i++ {
		dstFieldType := dstType.Field(i)
		dstFieldValue := dstValue.Field(i)
		if !dstFieldValue.CanSet() {
			return fmt.Errorf("field %s con't set", dstFieldType.Name)
		}

		tag := dstFieldType.Tag

		if tag.Get(tagIgnore) == "true" {
			continue
		}

		rule := ""
		if ruleTag := tag.Get(tagRule); ruleTag != "" {
			rule = ruleTag
		} else {
			rule = ruleDefault
		}

		var src []string
		if srcTag := tag.Get(tagSrc); srcTag != "" {
			src = strings.Split(srcTag, ",")
		} else {
			src = []string{dstFieldType.Name}
		}

		in, err := getValues(values, src)
		if err != nil {
			return err
		}

		err = setValue(in, dstFieldValue, rule)
		if err != nil {
			return err
		}
	}

	return nil
}

func getValues(values map[string]reflect.Value, src []string) ([]reflect.Value, error) {
	srcValues := make([]reflect.Value, 0, len(src))
	for _, s := range src {
		// 判断s的类型
		bs := []byte(s)
		lastIndex := len(bs) - 1
		if bs[lastIndex] == ']' {
			// slice 格式: Name[1]
			// 搜索']'，解析格式
			for i := lastIndex; i > 0; i-- {
				if bs[i] == '[' {
					vi, err := strconv.ParseInt(string(bs[i+1:lastIndex]), 10, 32)
					if err != nil {
						return nil, fmt.Errorf("src %s:%s format invalid", s, string(bs[i+1:lastIndex+1]))
					}

					index := int(vi)

					name := string(bs[0:i])
					v, ok := values[name]
					if !ok {
						return nil, fmt.Errorf("not found src field %s", name)
					}

					if index > v.Len()-1 {
						return nil, fmt.Errorf("index out of range by src field %s", name)
					}

					srcValues = append(srcValues, v.Index(index))
				}
			}
		} else {
			// default
			v, ok := values[s]
			if !ok {
				return nil, fmt.Errorf("not found src field %s", s)
			}

			srcValues = append(srcValues, v)
		}
	}

	return srcValues, nil
}

func setValue(in []reflect.Value, dst reflect.Value, rule string) error {
	if rule == ruleDefault {
		if len(in) != 1 {
			return errors.New("not set rule")
		}

		if !dst.Type().AssignableTo(in[0].Type()) {
			return fmt.Errorf("src %s can't assignable to dst %s", in[0].Type(), dst.Type())
		}

		dst.Set(in[0])

		return nil
	}

	ruleFunc, ok := getRule(rule)
	if !ok {
		return fmt.Errorf("not found rule %s", rule)
	}

	out := reflect.ValueOf(ruleFunc).Call(in)
	if len(out) == 2 {
		err, ok := out[1].Interface().(error)
		if ok && err != nil {
			return err
		}
	}

	if !dst.Type().AssignableTo(out[0].Type()) {
		return fmt.Errorf("src %s can't assignable to dst %s", out[0].Type().Name(), dst.Type().Name())
	}

	// dst.Type().
	dst.Set(out[0])

	return nil
}
