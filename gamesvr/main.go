package main

import (
	"gamesvr/controller"
	"gamesvr/manager"
	"gamesvr/session"
	"shared/protobuf/pb"
	"shared/utility/glog"
)

func main() {
	err := manager.Init(&session.Builder{})
	if err != nil {
		glog.Fatalf("manager init error: %v", err)
	}
	defer manager.ReporterClose()
	defer manager.LuaState.Close()

	// 监听信号
	go listenSignal()

	// 监听RPC
	handler, err := controller.NewGameHandler()
	if err != nil {
		glog.Fatalf("new handler error: %v", err)
	}
	handler.Init()

	pb.RegisterGameServer(manager.RPCServer, handler)
	err = manager.RPCServer.Listen("")
	if err != nil {
		glog.Fatalf("register and listen rpc error: %v", err)
	}
}
