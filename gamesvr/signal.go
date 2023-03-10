package main

import (
	"os"
	"os/signal"
	"syscall"

	"gamesvr/manager"
	"shared/utility/glog"
	"shared/utility/safe"
)

func listenSignal() {
	defer safe.Recover()

	ch := make(chan os.Signal, 1)

	signal.Notify(ch, syscall.SIGHUP, syscall.SIGABRT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGINT /*, syscall.SIGUSR1, syscall.SIGUSR2*/)

	for {
		s := <-ch
		switch s {

		case syscall.SIGABRT:
		case syscall.SIGHUP: // reload config
		case syscall.SIGINT, syscall.SIGTERM: // server close
			manager.SessManager.Close()

			err := manager.RPCServer.Stop()
			if err != nil {
				glog.Fatalf("failed to listen: %v", err)
			}

			// 保存至db
			// for _, session_ := range session.GetActiveSession() {
			// 	session_.Close()
			// }

			os.Exit(0)
		case syscall.SIGKILL:
		}
	}
}
