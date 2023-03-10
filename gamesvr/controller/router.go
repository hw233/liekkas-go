package controller

import (
	"reflect"

	"google.golang.org/protobuf/proto"

	"gamesvr/manager"
	"gamesvr/session"
	"shared/utility/router"
)

func NewGameRoute() (*router.Router, error) {
	gsRouter := router.NewRouter()

	err := gsRouter.RegisterHandler(&session.Session{}, router.WithConfig(manager.CSV.Protocol.Protocols()), router.WithFilter(filter))
	if err != nil {
		return nil, err
	}

	return gsRouter, nil
}

func filter(f interface{}) bool {
	t := reflect.TypeOf(f)

	// check kind
	if t.Kind() != reflect.Func {
		return false
	}

	// check in
	if t.NumIn() != 3 {
		return false
	}

	if t.In(0).String() != "*session.Session" {
		return false
	}

	if t.In(1).String() != "context.Context" {
		return false
	}

	_, ok := reflect.New(t.In(2)).Elem().Interface().(proto.Message)
	if !ok {
		return false
	}

	// check out
	if t.NumOut() != 2 {
		return false
	}

	_, ok = reflect.New(t.Out(0)).Elem().Interface().(proto.Message)
	if !ok {
		return false
	}

	if t.Out(1).String() != "error" {
		return false
	}

	return true
}
