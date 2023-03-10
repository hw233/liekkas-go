package module

import (
	"context"

	"github.com/go-redis/redis/v8"
	clientv3 "go.etcd.io/etcd/client/v3"

	"shared/utility/key"
)

// server discover
type Discover struct {
	prefix string
	client *clientv3.Client
}

func NewDiscover(service string, client *clientv3.Client) *Discover {
	return &Discover{
		prefix: key.MakeEtcdKey(service),
		client: client,
	}
}

func (d *Discover) SetServer(ctx context.Context, server, addr string) error {
	_, err := d.client.Put(ctx, d.makeKey(server), addr)
	if err != nil {
		return err
	}

	return nil
}

func (d *Discover) DelServer(ctx context.Context, server string) error {
	_, err := d.client.Delete(ctx, d.makeKey(server))
	if err != nil {
		return err
	}

	return nil
}

func (d *Discover) ListAddrs(ctx context.Context) ([]string, error) {
	resp, err := d.client.Get(ctx, d.prefix, clientv3.WithPrefix())
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, err
	}

	ret := make([]string, 0, len(resp.Kvs))

	for _, v := range resp.Kvs {
		ret = append(ret, string(v.Value))
	}

	return ret, nil
}

func (d *Discover) ListServers(ctx context.Context) ([]string, error) {
	resp, err := d.client.Get(ctx, d.prefix, clientv3.WithPrefix())
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, err
	}

	ret := make([]string, 0, len(resp.Kvs))

	for _, v := range resp.Kvs {
		keys := key.SplitEtcdKey(string(v.Key))
		if len(keys) >= 2 {
			ret = append(ret, keys[1])
		}
	}

	return ret, nil
}

func (d *Discover) makeKey(server string) string {
	return key.MakeEtcdKey(d.prefix, server)
}
