package etcd

//
// type Client struct {
// 	*clientv3.Client
// }
//
// func (c *Client) DialGRPC(service string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
// 	etcdResolver, err := resolver.NewBuilder(c.Client)
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return grpc.Dial("etcd:///"+service, append(opts, grpc.WithResolvers(etcdResolver))...)
// }
//
// func (c *Client) Register(service, addr string) error {
// 	em, err := endpoints.NewManager(c.Client, service)
// 	if err != nil {
// 		return err
// 	}
//
// 	return em.AddEndpoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Target: addr, Metadata: 0})
// }
//
// func (c *Client) UpdateSessionNum(service, addr string, num int) error {
// 	em, err := endpoints.NewManager(c.Client, service)
// 	if err != nil {
// 		return err
// 	}
//
// 	return em.AddEndpoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Target: addr, Metadata: num})
// }
//
// func (c *Client) Delete(service, addr string) error {
// 	em, err := endpoints.NewManager(c.Client, service)
// 	if err != nil {
// 		return err
// 	}
// 	return em.DeleteEndpoint(c.Ctx(), service+"/"+addr)
// }
//
// func (c *Client) RegisterWithLease(service, addr string, lid clientv3.LeaseID) error {
// 	em, err := endpoints.NewManager(c.Client, service)
// 	if err != nil {
// 		return err
// 	}
// 	return em.AddEndpoint(c.Ctx(), service+"/"+addr, endpoints.Endpoint{Target: addr}, clientv3.WithLease(lid))
// }
//
// func (c *Client) SaveConnStatus(service, id, addr string) error {
// 	// conn.status:portal:2100 127.0.0.1:8090
// 	// conn.status:gamesvr:2100 127.0.0.1:8080
//
// 	_, err := c.Put(c.Ctx(), fmt.Sprintf("conn.status:%s:%s", service, id), addr)
// 	return err
// }
//
// func (c *Client) GetConnStatus(service, id string) (string, error) {
// 	// conn.status:portal:2100 127.0.0.1:8090
// 	// conn.status:gamesvr:2100 127.0.0.1:8080
//
// 	resp, err := c.Get(c.Ctx(), fmt.Sprintf("conn.status:%s:%s", service, id))
// 	if err != nil {
// 		return "", err
// 	}
//
// 	if resp == nil || len(resp.Kvs) == 0 {
// 		return "", nil
// 	}
//
// 	return string(resp.Kvs[0].Value), nil
// }
