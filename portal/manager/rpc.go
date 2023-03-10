package manager

import (
	"shared/common"
	"shared/grpc"
	"shared/protobuf/pb"
	"shared/utility/glog"

	"github.com/openzipkin/zipkin-go/reporter"
)

var (
	RPCClient *grpc.Client
	RPCServer *grpc.Server

	RPCGameClient            pb.GameClient
	ZipkinReporterClientGame reporter.Reporter
)

func initRPC() error {
	var err error

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)
	RPCServer = grpc.NewServer(Conf.Service, Conf.ServerName, EtcdClient, RedisClient)

	// ---------------------------------------------------
	gameClientConn, r, err := RPCClient.Dial(common.ServiceGame, Conf.ServerName)
	ZipkinReporterClientGame = r
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	gameClientConn.RegisterFetcherFunc(func(args interface{}) int64 {
		if f, ok := args.(interface {
			GetUid() int64
		}); ok {
			return f.GetUid()
		}

		return 0
	})

	RPCGameClient = pb.NewGameClient(gameClientConn)

	// ---------------------------------------------------

	return nil
}

func ReporterClose() {
	ZipkinReporterClientGame.Close()
}
