package balancer

import (
	"context"
)

const (
	TypeRand  = 0
	TypeRound = 1
	TypeLeast = 2
)

type metadataKey struct{}

type Metadata struct {
	LBStrategy int
	Service    string
	Server     string
}

func GetCtxMetadata(ctx context.Context) (*Metadata, bool) {
	md, ok := ctx.Value(metadataKey{}).(*Metadata)
	return md, ok
}

func WithCtxMetadata(ctx context.Context, server string) context.Context {
	return context.WithValue(ctx, metadataKey{}, &Metadata{Server: server})
}

type Context struct {
	Server string
	ID     int64
}

func (c *Context) GetServer() (string, bool) {
	if c.Server == "" {
		return c.Server, false
	}

	return c.Server, true
}

func (c *Context) GetID() (int64, bool) {
	if c.ID == 0 {
		return c.ID, false
	}

	return c.ID, true
}

func WithContext(ctx context.Context, service string, c *Context) context.Context {
	return context.WithValue(ctx, service, c)
}

func GetContext(ctx context.Context, service string) (*Context, bool) {
	c, ok := ctx.Value(service).(*Context)
	if !ok {
		return &Context{}, false
	}

	return c, true
}
