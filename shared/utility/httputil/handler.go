package httputil

type HttpHandler interface {
	GetMethod() string
	FilterMethod(f interface{}) bool
}

type HttpMethodAdapter struct {
	Method string
}

func NewHttpMethodAdapter(Method string) *HttpMethodAdapter {
	return &HttpMethodAdapter{
		Method: Method,
	}
}

func (h *HttpMethodAdapter) GetMethod() string {
	return h.Method
}
