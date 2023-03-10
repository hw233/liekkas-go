package main

import (
	"net/http"

	"login/controller"
	"login/internal/base"
	"login/service"
	"shared/protobuf/pb"
	"shared/utility/glog"
)

func main() {
	base.Init()

	pb.RegisterLoginServer(base.RPCServer, service.RPCLoginHandler{})
	go base.RPCServer.Listen(base.Config.GRPCListenPort)

	glog.Infof("start listen http port: %s", base.Config.HTTPListenPort)
	if err := http.ListenAndServe("0.0.0.0:"+base.Config.HTTPListenPort, httpHandler()); err != nil {
		glog.Fatal("http.ListenAndServe error: ", err)
	}
}

func httpHandler() http.Handler {
	handler := http.NewServeMux()
	handler.HandleFunc("/v1/inner-register", controller.InnerRegister)
	handler.HandleFunc("/v1/inner-login", controller.InnerLogin)
	handler.HandleFunc("/v1/third-party-login", controller.ThirdPartyLogin)
	return handler
}
