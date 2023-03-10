package event

import (
	"shared/utility/param"
	"sync"
)

var (
	UserEventHandler = NewUserEventHandlers()
)

type UserEventHandlers struct {
	sync.RWMutex
	handlers map[int64]*EventHandlers
}

func NewUserEventHandlers() *UserEventHandlers {

	return &UserEventHandlers{
		handlers: map[int64]*EventHandlers{},
	}
}

func (e *UserEventHandlers) GetHandler(userId int64, Type int) (EventHandler, bool) {
	e.Lock()
	defer e.Unlock()
	v, ok := e.handlers[userId]
	if !ok {
		return nil, false
	}
	return v.GetHandler(Type)
}
func (e *UserEventHandlers) Register(userId int64, Type int, handler EventHandler) {
	e.Lock()
	defer e.Unlock()
	v, ok := e.handlers[userId]
	if !ok {
		v = NewEventHandlers()
		e.handlers[userId] = v
	}
	v.Register(Type, handler)
}

type EventHandlers map[int]EventHandler

func NewEventHandlers() *EventHandlers {
	return (*EventHandlers)(&map[int]EventHandler{})
}
func (e *EventHandlers) GetHandler(Type int) (EventHandler, bool) {
	v, ok := (*e)[Type]
	return v, ok
}
func (e *EventHandlers) Register(Type int, handler EventHandler) {
	(*e)[Type] = handler
}

type EventHandler func(Param *param.Param) error
