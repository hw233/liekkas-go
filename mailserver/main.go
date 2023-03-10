package main

import (
	"mailserver/controller"
	"mailserver/mail"
	"mailserver/manager"
	"shared/protobuf/pb"
	"shared/utility/glog"
)

func main() {
	err := manager.Init()
	if err != nil {
		glog.Fatalf("manager init error: %v", err)
		return
	}

	defer manager.ReporterClose()

	err = mail.Init()
	if err != nil {
		glog.Fatalf("mail manager init error: %v", err)
		return
	}

	serviceHandler := controller.NewServiceHandler()
	pb.RegisterMailServer(manager.RPCServer, serviceHandler)
	err = manager.RPCServer.Listen("")
	if err != nil {
		glog.Fatalf("register and listen rpc error: %v", err)
		return
	}
}
