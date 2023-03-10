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

	RPCGameClient     pb.GameClient
	RPCMailClient     pb.MailClient
	RPCForeplayClient pb.ForeplayClient

	ZipkinReporterClientGame     reporter.Reporter
	ZipkinReporterClientMail     reporter.Reporter
	ZipkinReporterClientForePlay reporter.Reporter
)

func initRPC() error {
	var err error

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)

	// ---------------------------------------------------
	gameClientConn, gameReporter, err := RPCClient.Dial(common.ServiceGame, "gm")
	ZipkinReporterClientGame = gameReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	RPCGameClient = pb.NewGameClient(gameClientConn)

	gameClientConn.RegisterBroadcastMethod("/Game/ReloadWhiteList")
	gameClientConn.RegisterBroadcastMethod("/Game/WhiteListSwitch")
	gameClientConn.RegisterBroadcastMethod("/Game/ReloadAnnouncement")
	gameClientConn.RegisterBroadcastMethod("/Game/ReloadMaintain")
	gameClientConn.IgnoreNullServer()
	// ---------------------------------------------------
	mailClientConn, mailReporter, err := RPCClient.Dial(common.ServiceMail, "gm")
	ZipkinReporterClientMail = mailReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	RPCMailClient = pb.NewMailClient(mailClientConn)

	foreplayClientConn, foreplayReporter, err := RPCClient.Dial(common.ServiceForeplay, "gm")
	ZipkinReporterClientForePlay = foreplayReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}
	RPCForeplayClient = pb.NewForeplayClient(foreplayClientConn)
	foreplayClientConn.RegisterBroadcastMethod("/Foreplay/ReloadAnnouncement")
	foreplayClientConn.RegisterBroadcastMethod("/Foreplay/ReloadMaintain")
	foreplayClientConn.IgnoreNullServer()

	return nil
}

func ReporterClose() {
	ZipkinReporterClientForePlay.Close()
	ZipkinReporterClientMail.Close()
	ZipkinReporterClientGame.Close()
}
