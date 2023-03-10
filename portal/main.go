package main

import (
	"encoding/binary"
	"fmt"
	"runtime"

	"portal/controller"
	"portal/manager"
	"shared/protobuf/pb"
	"shared/utility/glog"

	"github.com/panjf2000/gnet"
	"go.uber.org/zap/zapcore"
)

func main() {
	err := manager.Init()
	if err != nil {
		glog.Fatalf("manager init error: %v", err)
	}
	defer manager.ReporterClose()

	handler := controller.NewHandler()

	pb.RegisterPortalServer(manager.RPCServer, handler)
	go manager.RPCServer.Listen("")

	err = gnet.Serve(
		handler,
		fmt.Sprintf("tcp://:%s", manager.Conf.TCPListenPort),
		// options
		gnet.WithMulticore(true),
		gnet.WithReusePort(true),
		gnet.WithTCPKeepAlive(0),
		gnet.WithLoadBalancing(gnet.LeastConnections),
		gnet.WithNumEventLoop(runtime.NumCPU()),
		gnet.WithLogLevel(zapcore.DebugLevel),
		gnet.WithTCPNoDelay(gnet.TCPNoDelay),
		// gnet.WithSocketRecvBuffer(),
		// gnet.WithSocketSendBuffer(),
		// 4内容长度+2服务ID+内容
		gnet.WithCodec(gnet.NewLengthFieldBasedFrameCodec(
			gnet.EncoderConfig{
				ByteOrder:         binary.BigEndian,
				LengthAdjustment:  0,
				LengthFieldLength: 4,
			},
			gnet.DecoderConfig{
				ByteOrder:           binary.BigEndian,
				LengthFieldOffset:   0,
				LengthFieldLength:   4,
				LengthAdjustment:    0,
				InitialBytesToStrip: 4, //  去掉长度
			},
		)),
	)

	if err != nil {
		glog.Fatal(err)
	}
}
