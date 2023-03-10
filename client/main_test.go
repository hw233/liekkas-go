package main

import "testing"

var Pool *ClientPool

func TestMain(m *testing.M) {
	Pool = NewClientPool("127.0.0.1:9080")
	// defer Pool.Close()
	m.Run()
}
