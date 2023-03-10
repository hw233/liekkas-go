package main

import (
	"fmt"
	"net/http"

	"foreplay/controller"
	"foreplay/manager"
	"shared/protobuf/pb"
	"shared/utility/glog"
)

func main() {
	err := manager.Init()
	if err != nil {
		glog.Fatalf("manager init error: %v", err)
	}

	controller.RegisterHttpHandler()

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", manager.Conf.HTTPListenPort),
		Handler: nil,
	}
	err = httpServer.ListenAndServe()
	if err != nil {
		glog.Fatal("http listen fail: ", err)
	}

	// 监听RPC
	rpcHandler, err := controller.NewForeplayHandler()
	if err != nil {
		glog.Fatalf("new handler error: %v", err)
	}

	pb.RegisterForeplayServer(manager.RPCServer, rpcHandler)
	err = manager.RPCServer.Listen("")
	if err != nil {
		glog.Fatalf("register and listen rpc error: %v", err)
	}
}
