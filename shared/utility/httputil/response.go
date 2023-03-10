package httputil

import (
	"shared/utility/errors"
)

type HttpContent struct {
	ErrorCode int         `json:"err_code"`
	ErrorMsg  string      `json:"message"` // gm平台接入规范
	Data      interface{} `json:"data"`
}

func NewHttpContent(data interface{}, err error) *HttpContent {
	if err == nil {
		return &HttpContent{
			Data: data,
		}
	} else {
		return &HttpContent{
			ErrorCode: errors.Code(err),
			ErrorMsg:  errors.Text(err),
		}
	}

}
