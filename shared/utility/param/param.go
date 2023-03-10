package param

import (
	"encoding/json"
	"errors"
	"strconv"
)

var ErrIndexOverParams = errors.New("index over params")

type Param struct {
	params []string
}

func NewParam(params []string) *Param {
	return &Param{
		params: params,
	}
}

func (p *Param) Len() int {
	return len(p.params)
}

func (p *Param) GetString(i int) (string, error) {
	if i >= len(p.params) {
		return "", ErrIndexOverParams
	}

	return p.params[i], nil
}

func (p *Param) GetInterface(i int) (interface{}, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	return interface{}(p.params[i]), nil
}

func (p *Param) GetInt(i int) (int, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseInt(p.params[i], 10, 64)
	if err != nil {
		return 0, err
	}

	return int(ret), nil
}

func (p *Param) GetInt8(i int) (int8, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseInt(p.params[i], 10, 8)
	if err != nil {
		return 0, err
	}

	return int8(ret), nil
}

func (p *Param) GetInt16(i int) (int16, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseInt(p.params[i], 10, 16)
	if err != nil {
		return 0, err
	}

	return int16(ret), nil
}

func (p *Param) GetInt32(i int) (int32, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseInt(p.params[i], 10, 32)
	if err != nil {
		return 0, err
	}

	return int32(ret), nil
}

func (p *Param) GetInt64(i int) (int64, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseInt(p.params[i], 10, 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (p *Param) GetUint8(i int) (uint8, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseUint(p.params[i], 10, 8)
	if err != nil {
		return 0, nil
	}

	return uint8(ret), nil
}

func (p *Param) GetUint16(i int) (uint16, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseUint(p.params[i], 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(ret), nil
}

func (p *Param) GetUint32(i int) (uint32, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseUint(p.params[i], 10, 32)
	if err != nil {
		return 0, err
	}

	return uint32(ret), nil
}

func (p *Param) GetUint64(i int) (uint64, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseUint(p.params[i], 10, 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (p *Param) GetFloat32(i int) (float32, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseFloat(p.params[i], 32)
	if err != nil {
		return 0, err
	}

	return float32(ret), nil
}

func (p *Param) GetFloat64(i int) (float64, error) {
	if i >= len(p.params) {
		return 0, ErrIndexOverParams
	}

	ret, err := strconv.ParseFloat(p.params[i], 64)
	if err != nil {
		return 0, err
	}

	return ret, nil
}

func (p *Param) GetBool(i int) (bool, error) {
	if i >= len(p.params) {
		return false, ErrIndexOverParams
	}

	ret, err := strconv.ParseBool(p.params[i])
	if err != nil {
		return false, err
	}

	return ret, nil
}

func (p *Param) GetBytes(i int) ([]byte, error) {
	if i >= len(p.params) {
		return []byte{}, ErrIndexOverParams
	}

	return []byte(p.params[i]), nil
}

func (p *Param) GetInt32s(i int) ([]int32, error) {
	bs, err := p.GetBytes(i)
	if err != nil {
		return []int32{}, err
	}

	var ret []int32

	err = json.Unmarshal(bs, &ret)
	if err != nil {
		return []int32{}, err
	}

	return ret, nil
}

func (p *Param) GetInt64s(i int) ([]int64, error) {
	bs, err := p.GetBytes(i)
	if err != nil {
		return []int64{}, err
	}

	var ret []int64

	err = json.Unmarshal(bs, &ret)
	if err != nil {
		return []int64{}, err
	}

	return ret, nil
}
