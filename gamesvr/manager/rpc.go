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

	RPCPortalClient pb.PortalClient
	RPCGuildClient  pb.GuildClient
	RPCMailClient   pb.MailClient
	RPCLoginClient  pb.LoginClient

	ZipkinReporterClientPortal reporter.Reporter
	ZipkinReporterClientGuild  reporter.Reporter
	ZipkinReporterClientMail   reporter.Reporter
	ZipkinReporterClientLogin  reporter.Reporter
)

func initRPC() error {
	var err error

	RPCClient = grpc.NewClient(EtcdClient, RedisClient)
	RPCServer = grpc.NewServer(Conf.Service, Conf.ServerName, EtcdClient, RedisClient)

	// ---------------------------------------------------
	portalClientConn, portalReporter, err := RPCClient.Dial(common.ServicePortal, Conf.Service)
	ZipkinReporterClientPortal = portalReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	portalClientConn.RegisterFetcherFunc(func(args interface{}) int64 {
		if f, ok := args.(interface {
			GetUid() int64
		}); ok {
			return f.GetUid()
		}

		return 0
	})

	RPCPortalClient = pb.NewPortalClient(portalClientConn)

	// ---------------------------------------------------

	guildClientConn, guildReporter, err := RPCClient.Dial(common.ServiceGuild, Conf.Service)
	ZipkinReporterClientGuild = guildReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	guildClientConn.RegisterFetcherFunc(func(args interface{}) int64 {
		if f, ok := args.(interface {
			GetGuildID() int64
		}); ok {
			return f.GetGuildID()
		}

		return 0
	})

	RPCGuildClient = pb.NewGuildClient(guildClientConn)

	// ---------------------------------------------------

	mailClientConn, mailReporter, err := RPCClient.Dial(common.ServiceMail, Conf.Service)
	ZipkinReporterClientMail = mailReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	RPCMailClient = pb.NewMailClient(mailClientConn)
	// ---------------------------------------------------

	loginClientConn, loginReporter, err := RPCClient.Dial(common.ServiceLogin, Conf.Service)
	ZipkinReporterClientLogin = loginReporter
	if err != nil {
		glog.Errorf("RPCClient.Dial() error: %v", err)
		return err
	}

	RPCLoginClient = pb.NewLoginClient(loginClientConn)

	return nil
}

func ReporterClose() {
	ZipkinReporterClientLogin.Close()
	ZipkinReporterClientMail.Close()
	ZipkinReporterClientGuild.Close()
	ZipkinReporterClientPortal.Close()
}
