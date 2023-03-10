package router

import (
	"testing"
)

type testHandler struct {
	i int
}

func (th *testHandler) Func1() {
	th.i += 1
}

func (th *testHandler) Func2() {
	th.i += 2
}

func (th *testHandler) Func3(i int) {
	th.i += i
}

func TestRouter_RegisterByConfigOption(t *testing.T) {
	config := map[int32]string{
		1: "Func1",
		3: "Func3",
	}

	handler := &testHandler{}
	router := NewRouter()
	err := router.RegisterHandler(handler, WithConfig(config))
	if err != nil {
		t.Fatalf("router.RegisterByHandler() error: %v", err)
	}

	f1, ok := router.Route(1)
	if !ok {
		t.Fatalf("router.GetHandleFunc(0) !ok")
	}

	f1.Interface().(func(th *testHandler))(handler)
	if handler.i != 1 {
		t.Fatalf("handler.i != 1")
	}

	f2, ok := router.Route(3)
	if !ok {
		t.Fatalf("router.GetHandleFunc(1) !ok")
	}

	f2.Interface().(func(th *testHandler, i int))(handler, 3)
	if handler.i != 4 {
		t.Fatalf("handler.i != 4")
	}
}

// func TestRouter_RegisterByFilterOption(t *testing.T) {
// 	config := map[int32]string{
// 		1: "Func1",
// 		2: "Func2",
// 		3: "Func3",
// 	}
//
// 	handler := &testHandler{}
// 	router := NewRouter()
// 	err := router.RegisterHandler(handler, WithConfig(config), WithFilter(func(th *testHandler) {}))
// 	if err != nil {
// 		t.Fatalf("router.RegisterByHandler() error: %v", err)
// 	}
//
// 	f1, ok := router.Route(1)
// 	if !ok {
// 		t.Fatalf("router.GetHandleFunc(%d) !ok", 1)
// 	}
//
// 	f1.Interface().(func(th *testHandler))(handler)
// 	if handler.i != 1 {
// 		t.Fatalf("handler.i != 1")
// 	}
//
// 	f2, ok := router.Route(2)
// 	if !ok {
// 		t.Fatalf("router.GetHandleFunc(%d) !ok", 2)
// 	}
//
// 	f2.Interface().(func(th *testHandler))(handler)
// 	if handler.i != 3 {
// 		t.Fatalf("handler.i != 3")
// 	}
//
// 	// filter func
// 	_, ok = router.Route(3)
// 	if ok {
// 		t.Fatalf("router.GetHandleFunc(%d) ok", 3)
// 	}
// }

func TestRouter_RouteCall(t *testing.T) {
	config := map[int32]string{
		1: "Func1",
		3: "Func3",
	}

	handler := &testHandler{}
	router := NewRouter()
	err := router.RegisterHandler(handler, WithConfig(config))
	if err != nil {
		t.Fatalf("router.RegisterByHandler() error: %v", err)
	}

	router.RouteCall(1, handler)
	if handler.i != 1 {
		t.Fatalf("handler.i != 1")
	}

	router.RouteCall(3, handler, 3)
	if handler.i != 4 {
		t.Fatalf("handler.i != 4")
	}
}

// BenchmarkRouter_Route-8       	23535302	        49.17 ns/op
func BenchmarkRouter_Route(b *testing.B) {
	handler := &testHandler{}
	router := NewRouter()

	// register 10000 function
	for i := 0; i < 10000; i++ {
		_ = router.Register(int32(i), func(th *testHandler) {})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		f1, _ := router.Route(int32(i % 9999))
		f1.Interface().(func(th *testHandler))(handler)
	}
}

// BenchmarkRouter_RouteCall-8   	 5471329	       218.2 ns/op
func BenchmarkRouter_RouteCall(b *testing.B) {
	handler := &testHandler{}
	router := NewRouter()

	// register 10000 function
	for i := 0; i < 10000; i++ {
		_ = router.Register(int32(i), func(th *testHandler) {})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		router.RouteCall(int32(i%9999), handler)
	}
}
