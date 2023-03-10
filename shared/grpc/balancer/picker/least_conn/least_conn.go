package least_conn

import (
	"bytes"
	"context"
	"fmt"
	"sync"

	"shared/grpc/balancer"
	"shared/grpc/module"
	"shared/utility/errors"
	"shared/utility/glog"

	"github.com/go-redis/redis/v8"
	ggbalancer "google.golang.org/grpc/balancer"
	ggbase "google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	pickBuilderBuildLog = `
--------- service [%s] regenerate server ----------
%s---------------------------------------------------`
)

// Name is the name of least_conn balancer.
const Name = "least_conn"

type lcPickerBuilder struct {
	client  *redis.Client
	service string
	pick    *picker
}

func NewLCPickerBuilder(client *redis.Client, service string) *lcPickerBuilder {
	return &lcPickerBuilder{
		client:  client,
		service: service,
	}
}

func (pb *lcPickerBuilder) Build(buildInfo ggbase.PickerBuildInfo) ggbalancer.Picker {
	conns := make(map[string]ggbalancer.SubConn, len(buildInfo.ReadySCs))

	connCount := 1
	buffer := bytes.Buffer{}

	for subConn, subConnInfo := range buildInfo.ReadySCs {
		conns[subConnInfo.Address.ServerName] = subConn

		buffer.WriteString(fmt.Sprintf("          %d. %s, %s\n", connCount, subConnInfo.Address.ServerName, subConnInfo.Address.Addr))
		connCount++
	}

	if connCount > 1 {
		glog.Infof(pickBuilderBuildLog, pb.service, buffer.String())
	}

	return &picker{
		service:  pb.service,
		conns:    conns,
		Balancer: module.NewBalancer(pb.service, pb.client),
		// recorder: module.NewRecorder(pb.client),
	}
}

type picker struct {
	sync.RWMutex
	service string
	conns   map[string]ggbalancer.SubConn
	*module.Balancer
	// recorder *module.Recorder
}

// func newPicker(client *redis.Client, service string) *picker {
// 	return &picker{
// 		service:  service,
// 		conns:    map[string]ggbalancer.SubConn{},
// 		balancer: module.NewBalancer(client),
// 		recorder: module.NewRecorder(client),
// 	}
// }

func (p *picker) Pick(info ggbalancer.PickInfo) (ggbalancer.PickResult, error) {
	p.RLock()
	defer p.RUnlock()

	ret := ggbalancer.PickResult{}

	subConn, err := p.pick(info.Ctx)
	if err != nil {
		glog.Errorf("picker: pick error: %v", err)
		return ret, err
	}

	ret.SubConn = subConn

	return ret, nil
}

func (p *picker) pick(ctx context.Context) (ggbalancer.SubConn, error) {
	md, ok := balancer.GetCtxMetadata(ctx)
	if !ok {
		return nil, errors.New("not found metadata")
	}

	// glog.Debugf("md: %+v", md)
	// c, ok := GetContext(ctx, p.service)
	// if !ok {
	// 	// TODO：先跑下去，后面报错，在picker做随机服务器处理
	// }

	// 自带server
	if md.Server != "" {
		subConn, ok := p.conns[md.Server]
		if !ok {
			// 目标服务器发生异常
			return nil, status.Errorf(codes.Unavailable, "target %s disable", md.Server)
		}

		return subConn, nil
	}

	// 从 recorder 取服务器
	// server, err := p.recorder.GetBalance(ctx, p.service, c.ID)
	// if err != nil {
	// 	return nil, err
	// }

	// 有了
	// if server != "" {
	// 	subConn, ok := p.conns.GetBalance(server)
	// 	if !ok {
	// 		// 目标服务器发生异常
	// 		glog.Debug("picker: conns: %v", p.conns)
	// 		return nil, status.Errorf(codes.Unavailable, "name %s disable", server)
	// 	}
	//
	// 	c.Server = server
	//
	// 	return subConn, nil
	// }

	// 从 balancer 取服务器
	server, err := p.GetBalance(ctx)
	if err != nil {
		glog.Errorf("pick: balancer.GetBalance(%s) error: %v", err)
		return nil, err
	}

	err = p.IncrBalance(ctx, server)
	if err != nil {
		glog.Errorf("pick: balancer.IncrBalance(%s) error: %v", err)
		return nil, err
	}

	subConn, ok := p.conns[server]
	if !ok {
		glog.Errorf("pick: conns.GetBalance(%s) !ok", server)
		// 目标服务器发生异常
		return nil, status.Errorf(codes.Unavailable, "name %s disable", server)
	}

	return subConn, nil
}
