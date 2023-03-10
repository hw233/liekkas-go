package grpc

import (
	"context"
	"fmt"
	"net"
	"strings"

	"shared/grpc/module"
	"shared/utility/glog"
	"shared/utility/ip"

	"github.com/go-redis/redis/v8"
	"github.com/openzipkin/zipkin-go/reporter"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"

	zipkin "github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

type Server struct {
	*grpc.Server
	registerService []string
	*module.Balancer
	*module.Discover
	*module.Recorder
	Service        string
	ServerName     string
	ZipkinReporter reporter.Reporter
}

func NewServer(service, serverName string, etcdClient *clientv3.Client, redisClient *redis.Client) *Server {
	tracer, r, err := NewZipkinTracer(service, "localhost:0")
	if err != nil {
		glog.Errorf("NewZipkinTracer Wrong")
	}

	return &Server{
		Server:          grpc.NewServer(grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer))),
		registerService: []string{},
		Balancer:        module.NewBalancer(service, redisClient),
		Discover:        module.NewDiscover(service, etcdClient),
		Recorder:        module.NewRecorder(service, redisClient),
		Service:         service,
		ServerName:      serverName,
		ZipkinReporter:  r,
	}
}

func NewZipkinTracer(serviceName, hostPort string) (*zipkin.Tracer, reporter.Reporter, error) {

	// 初始化zipkin reporter
	// reporter可以有很多种，如：logReporter、httpReporter
	r := httpreporter.NewReporter("http://localhost:9411/api/v2/spans")
	// r := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	// r := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))

	//创建一个endpoint，用来标识当前服务，服务名：服务地址和端口
	endpoint, err := zipkin.NewEndpoint(serviceName, hostPort)
	if err != nil {
		return nil, r, err
	}

	// 初始化追踪器 主要作用有解析span，解析上下文等
	tracer, err := zipkin.NewTracer(r, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		return nil, r, err
	}

	return tracer, r, nil
}

func (s *Server) RegisterService(desc *grpc.ServiceDesc, impl interface{}) {
	// s.desc = desc
	// s.svr = impl
	// for i, method := range desc.Methods {
	// 	s.methods[method.MethodName] = desc.Methods[i]
	// }
	s.registerService = append(s.registerService, strings.ToLower(desc.ServiceName))
	s.Server.RegisterService(desc, impl)
}

// ex: ":9090"
func (s *Server) Listen(addr string) error {
	if addr == "" {
		localAddr, err := ip.LocalAddr()
		if err != nil {
			return err
		}

		// random port
		addr = fmt.Sprintf("%s:0", localAddr)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// clear count in balancer
	err = s.SetBalance(ctx, 0)
	if err != nil {
		return err
	}

	err = s.SetServer(ctx, lis.Addr().String())
	if err != nil {
		return err
	}

	glog.Infof("grpc listen addr: %s", lis.Addr().String())

	return s.Serve(lis)
}

func (s *Server) Stop() error {
	ctx := context.Background()

	// delete from discover
	err := s.DelServer(ctx)
	if err != nil {
		return err
	}

	// delete from balancer
	err = s.DelBalance(ctx)
	if err != nil {
		return err
	}

	s.Server.GracefulStop()
	s.ZipkinReporter.Close()

	return nil
}

// auto input serverName
func (s *Server) IncrBalance(ctx context.Context) error {
	return s.Balancer.IncrBalance(ctx, s.ServerName)
}

func (s *Server) DecrBalance(ctx context.Context) error {
	return s.Balancer.DecrBalance(ctx, s.ServerName)
}

func (s *Server) SetBalance(ctx context.Context, count int) error {
	return s.Balancer.SetBalance(ctx, s.ServerName, count)
}

func (s *Server) DelBalance(ctx context.Context) error {
	return s.Balancer.DelBalance(ctx, s.ServerName)
}

func (s *Server) SetServer(ctx context.Context, addr string) error {
	return s.Discover.SetServer(ctx, s.ServerName, addr)
}

func (s *Server) DelServer(ctx context.Context) error {
	return s.Discover.DelServer(ctx, s.ServerName)
}

func (s *Server) SetRecord(ctx context.Context, id int64) error {
	return s.Recorder.SetRecord(ctx, s.ServerName, id)
}
