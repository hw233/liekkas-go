package httputil

import (
	"net/url"
	"shared/utility/errors"
	"shared/utility/glog"
	"sort"
	"strconv"
	"strings"
)

type Params struct {
	values map[string]*Param
}

func NewParams(m url.Values) *Params {
	values := map[string]*Param{}
	for k, v := range m {
		values[k] = NewParam(k, v[0])
	}
	return &Params{
		values: values,
	}
}

func EmptyParams() *Params {
	return &Params{
		values: map[string]*Param{},
	}
}
func (p *Params) Put(k, v string) {
	p.values[k] = NewParam(k, v)
}

func (p *Params) Get(k string) *Param {
	param, ok := p.values[k]
	if !ok {
		return NewParam(k, "")
	}
	return param
}

// Encode 参数encode成 aa=bb&cc=dd的格式
func (p *Params) Encode() string {
	if len(p.values) == 0 {
		return ""
	}
	var buf strings.Builder
	keys := make([]string, 0, len(p.values))
	for k := range p.values {
		keys = append(keys, k)
	}
	// 一定要排序
	sort.Strings(keys)
	for _, k := range keys {
		vs := p.values[k]
		keyEscaped := url.QueryEscape(k)
		if buf.Len() > 0 {
			buf.WriteByte('&')
		}
		buf.WriteString(keyEscaped)
		buf.WriteByte('=')
		buf.WriteString(url.QueryEscape(vs.String()))

	}
	return buf.String()
}

func (p *Params) Values() map[string]*Param {
	return p.values
}

type Param struct {
	K, V string
}

func NewParam(k, v string) *Param {
	return &Param{
		K: k,
		V: v,
	}
}
func (p *Param) wrapErr(err error) error {
	glog.Debugf("parse error value err:%v", errors.WrapTrace(err))
	return errors.New("parse param: \"%s\" error value:\"%s\"", p.K, p.V)

}

func (p *Param) paramEmptyErr() error {
	return errors.New("param: \"%s\" is empty ", p.K)

}

func (p *Param) Bool() (bool, error) {
	if p.IsEmpty() {
		return false, p.paramEmptyErr()
	}

	ret, err := strconv.ParseBool(p.String())
	if err != nil {
		return false, p.wrapErr(err)
	}

	return ret, nil
}

func (p *Param) String() string {
	return p.V
}
func (p *Param) Int() (int, error) {
	if p.IsEmpty() {
		return 0, p.paramEmptyErr()
	}

	ret, err := strconv.ParseInt(p.String(), 10, 64)
	if err != nil {
		return 0, p.wrapErr(err)
	}

	return int(ret), nil
}

func (p *Param) Int32() (int32, error) {
	if p.IsEmpty() {
		return 0, p.paramEmptyErr()
	}

	ret, err := strconv.ParseInt(p.String(), 10, 32)
	if err != nil {
		return 0, p.wrapErr(err)
	}

	return int32(ret), nil
}

func (p *Param) Int64() (int64, error) {
	if p.IsEmpty() {
		return 0, p.paramEmptyErr()
	}

	ret, err := strconv.ParseInt(p.String(), 10, 64)
	if err != nil {
		return 0, p.wrapErr(err)
	}

	return ret, nil
}

func (p *Param) Int64s() ([]int64, error) {

	split := strings.Split(p.String(), ",")

	var ret []int64
	for _, s := range split {
		if s == "" {
			continue
		}
		num, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return nil, p.wrapErr(err)
		}
		ret = append(ret, num)
	}

	return ret, nil
}

func (p *Param) Int32s() ([]int32, error) {

	split := strings.Split(p.String(), ",")

	var ret []int32
	for _, s := range split {
		if s == "" {
			continue
		}

		num, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return nil, p.wrapErr(err)
		}
		ret = append(ret, int32(num))
	}

	return ret, nil
}

func (p *Param) IsEmpty() bool {
	return p.V == ""
}
