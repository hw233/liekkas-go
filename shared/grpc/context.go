package grpc

import (
	"context"

	"shared/utility/glog"
)

type metadataKey struct{}

type metadataMap struct {
	md map[string]*Metadata
}

type Metadata struct {
	servers []string
	err     error
}

func (m *Metadata) Server() string {
	if len(m.servers) > 0 {
		return m.servers[0]
	}

	return ""
}

func (m *Metadata) SetServer(server string) {
	if len(m.servers) > 1 {
		m.servers[0] = server
		return
	}

	m.servers = append(m.servers, server)
}

func (m *Metadata) Servers() []string {
	return m.servers
}

func (m *Metadata) Err() error {
	return m.err
}

func WithCtxErr(ctx context.Context, service string, err error) context.Context {
	mdm, ok := ctx.Value(metadataKey{}).(*metadataMap)
	if ok {
		mdm.md[service] = &Metadata{
			err: err,
		}

		return ctx
	}

	mdm = &metadataMap{md: map[string]*Metadata{}}
	mdm.md[service] = &Metadata{
		err: err,
	}

	return context.WithValue(ctx, metadataKey{}, mdm)
}

func WithCtxMetadata(ctx context.Context, service string, servers []string) context.Context {
	glog.Debugf("WithCtxMetadata: Service: %s, servers: %v", service, servers)
	mdm, ok := ctx.Value(metadataKey{}).(*metadataMap)
	if ok {
		mdm.md[service] = &Metadata{
			servers: servers,
		}

		return ctx
	}

	mdm = &metadataMap{md: map[string]*Metadata{}}
	mdm.md[service] = &Metadata{
		servers: servers,
	}

	return context.WithValue(ctx, metadataKey{}, mdm)
}

func GetCtxMetadata(ctx context.Context, service string) (*Metadata, bool) {
	mdm, ok := ctx.Value(metadataKey{}).(*metadataMap)
	if !ok {
		return nil, false
	}

	md, ok := mdm.md[service]
	if !ok {
		return nil, false
	}

	return md, true
}
