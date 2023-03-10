package main

import (
	"log"

	"guild/controller"
	"guild/manager"
	"guild/session"
	"shared/protobuf/pb"
)

func main() {
	err := manager.Init(&session.Builder{})
	if err != nil {
		log.Fatalf("manager init error: %v", err)
	}
	defer manager.ReporterClose()

	// 监听信号
	go listenSignal()

	// 监听RPC
	handler := controller.NewGuildHandler()

	pb.RegisterGuildServer(manager.RPCServer, handler)
	err = manager.RPCServer.Listen("")
	if err != nil {
		log.Fatalf("register and listen rpc error: %v", err)
	}
}
