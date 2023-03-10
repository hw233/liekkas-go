package event

import "fmt"

type Event struct {
	Type   int      `json:"type"`
	Params []string `json:"params"`
}

func NewEvent(Type int, params ...interface{}) *Event {
	Params := make([]string, 0, len(params))

	for _, param := range params {
		Params = append(Params, fmt.Sprint(param))
	}
	return &Event{
		Type:   Type,
		Params: Params,
	}
}
