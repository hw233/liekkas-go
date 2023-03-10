package main

import (
	"fmt"
	"gm/controller"
	"gm/manager"
	"net/http"
	"shared/utility/glog"
)

func main() {
	err := manager.Init()
	if err != nil {
		glog.Fatalf("manager init error: %v", err)
	}
	defer manager.ReporterClose()
	mux := http.NewServeMux()
	dispatcher := controller.NewDispatcher()
	err = dispatcher.RegisterRouters()
	if err != nil {
		glog.Fatalf("dispatcher RegisterAll error: %v", err)
	}
	dispatcher.AddFilter(&controller.BearerAuthenticationFilter{})
	mux.Handle("/", dispatcher)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", manager.Conf.Http.Port),
		Handler: mux,
	}

	glog.Info("Starting HTTP server ...")
	err = server.ListenAndServe()
	if err != nil {
		if err == http.ErrServerClosed {
			glog.Info("Server closed under request")
		} else {
			glog.Fatal("Server closed unexpected")
		}
	}
}
