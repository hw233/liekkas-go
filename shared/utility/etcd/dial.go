package etcd

import (
	"context"
	"shared/utility/glog"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc"
)

func Dial(endpoints []string) (*clientv3.Client, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
		DialOptions: []grpc.DialOption{grpc.WithBlock()}, // https://github.com/etcd-io/etcd/issues/9877
	})

	if err != nil {
		return nil, err
	}

	for _, v := range endpoints {
		_, err = client.Status(context.Background(), v)
		if err != nil {
			glog.Errorf("ERROR: etcd status error: %s:%v \n", v, err)
			return nil, err
		}

		glog.Infof("etcd %v connect success", v)
	}

	return client, nil
}
