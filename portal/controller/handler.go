package controller

import (
	"context"
	"google.golang.org/grpc/status"
	"time"

	"github.com/panjf2000/gnet"
	"golang.org/x/time/rate"

	"portal/base"
	"portal/manager"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type Handler struct {
	*gnet.EventServer
	*pb.UnimplementedPortalServer
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) OnOpened(c gnet.Conn) (out []byte, action gnet.Action) {
	c.SetContext(base.NewConn(c, time.Second*10, rate.NewLimiter(rate.Limit(manager.Conf.RequestRateSec), manager.Conf.RequestRateCap)))
	return
}

func (h *Handler) OnClosed(c gnet.Conn, err error) (action gnet.Action) {
	cc, ok := c.Context().(*base.Conn)
	if ok {
		manager.ConnPool.DelConn(cc)
	}

	err = manager.RPCServer.DelRecord(context.Background(), cc.UID)
	if err != nil {
		// handle error
	}

	return
}

func (h *Handler) OnInitComplete(srv gnet.Server) (action gnet.Action) {
	glog.Infof("Echo server is listening on %s (multi-cores: %t, loops: %d)\n",
		srv.Addr.String(), srv.Multicore, srv.NumEventLoop)
	return
}

func (h *Handler) React(frame []byte, c gnet.Conn) (out []byte, action gnet.Action) {
	cc, ok := c.Context().(*base.Conn)
	if !ok {
		glog.Errorf("c.Context().(*tcp.Conn) !ok")
		// 连接有问题，主动断掉
		_ = c.Close()
		return
	}

	// 异步处理
	// TODO：暂时这么写着，后面要控制每个玩家携程数，错误数，流量，频率，防止某些玩家异常行为和攻击
	if err := manager.GoPool.Submit(func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()

		// ctx = manager.RPCGameClient.WithSinglecastCtx(ctx, grpc.SinglecastOpt(cc.UID, cc.Server))

		err := h.HandleClient(ctx, cc, frame)
		if err != nil {
			glog.Errorf("handle conn content error: %v", err)
			return
		}
	}); err != nil {
		// 池子满了，丢掉数据
		glog.Errorf("pools full error: %v", err)
		return
	}

	return
}

func (h *Handler) Push(ctx context.Context, req *pb.PushReq) (*pb.PushResp, error) {
	glog.Debugf("Push: req[%+v]", req)
	if req == nil || req.GameResp == nil {
		return nil, errors.New("nil req")
	}

	c, ok := manager.ConnPool.GetConn(req.Uid)
	if !ok {
		return nil, errors.New("not found c")
	}

	out, err := proto.Marshal(req.GameResp)
	if err != nil {
		return nil, err
	}

	err = c.AsyncWrite(out)
	if err != nil {
		return nil, err
	}

	return &pb.PushResp{}, nil
}

func (h *Handler) HandleClient(ctx context.Context, conn *base.Conn, content []byte) error {
	gameReq := &pb.GameReq{}
	err := proto.Unmarshal(content, gameReq)
	if err != nil {
		return err
	}

	gameResp, err := h.handle(ctx, conn, gameReq)
	if err != nil {
		gameResp, err = h.Error(conn, err, gameReq)
		if err != nil {
			return err
		}
	}

	if gameResp != nil {
		out, err := proto.Marshal(gameResp)
		if err != nil {
			return err
		}

		// 发送给客户端
		err = conn.AsyncWrite(out)
		if err != nil {
			glog.Errorf("ERROR: AsyncWrite Error: %v", err)
			return err
		}
	}

	glog.Debugf("HandleClient: resp[%+v]", gameResp)

	return nil
}

func (h *Handler) handle(ctx context.Context, conn *base.Conn, gameReq *pb.GameReq) (*pb.GameResp, error) {
	if !conn.IsConnected {
		// 第一次连接

		// 检查最大连接数上限
		if manager.ConnPool.Len() >= manager.Conf.MaxConn {
			return nil, common.ErrTooManyConn
		}

		gameResp, err := h.Connect(ctx, conn, gameReq)
		if err != nil {
			return nil, err
		}

		// 连接成功加入连接池
		if conn.IsConnected {
			conn.SetContext(conn)
			manager.ConnPool.PutConn(conn)

			err = manager.RPCServer.SetRecord(ctx, conn.UID)
			if err != nil {
				return nil, err
			}
		}

		return gameResp, nil
	}

	gameResp, err := h.Transfer(ctx, conn, gameReq)
	if err != nil {
		return nil, err
	}

	return gameResp, nil
}

// 连接的目的是在portal拿到id
func (h *Handler) Connect(ctx context.Context, conn *base.Conn, gameReq *pb.GameReq) (*pb.GameResp, error) {
	conn.Lock()
	defer conn.Unlock()

	glog.Debugf("Connect: req[%+v]", gameReq)

	// TODO： 传conn的token过去验证conn的唯一性
	req := &pb.ConnectReq{
		GameReq: gameReq,
		Token:   conn.Token,
		Server:  manager.Conf.ServerName,
		Ctime:   conn.CTime,
	}

	resp, err := manager.RPCGameClient.Connect(ctx, req)
	if err != nil {
		return nil, err
	}

	glog.Debugf("Connect: resp[%s]", prototext.Format(resp))

	if resp.Uid != 0 {
		conn.UID = resp.Uid
		conn.IsConnected = true
	}

	return resp.GameResp, nil
}

func (h *Handler) redo(ctx context.Context, conn *base.Conn, clientReq *pb.GameReq) (*pb.GameResp, error) {
	conn.RLock()
	defer conn.RUnlock()
	glog.Debugf("Transfer: req[%+v]", clientReq)

	// 控制请求频率
	if !conn.Allow() {
		return nil, common.ErrTooManyRequest
	}

	req := &pb.TransferReq{
		GameReq: clientReq,
		Uid:     conn.UID,
		Token:   conn.Token,
		Ctime:   conn.CTime,
	}

	resp, err := manager.RPCGameClient.Transfer(ctx, req)
	if err != nil {
		return nil, err
	}

	glog.Debugf("redo: resp[%s]", prototext.Format(resp))

	// TODO：这里获取server的方式可以改变下，放到Connect获取
	if conn.Server == "" {
		// md, ok := manager.RPCGameClient.GetCtxMetadata(ctx)
		// if ok {
		// 	conn.Server = md.Server()
		// }
	}

	return resp.GameResp, nil

}

func (h *Handler) Transfer(ctx context.Context, conn *base.Conn, clientReq *pb.GameReq) (*pb.GameResp, error) {
	conn.RLock()
	defer conn.RUnlock()

	glog.Debugf("Transfer: req[%+v]", clientReq)

	// 控制请求频率
	if !conn.Allow() {
		return nil, common.ErrTooManyRequest
	}

	req := &pb.TransferReq{
		GameReq: clientReq,
		Uid:     conn.UID,
		Token:   conn.Token,
		Ctime:   conn.CTime,
	}

	resp, err := manager.RPCGameClient.Transfer(ctx, req)
	if err != nil {
		if _, ok := status.FromError(err); ok {
			//if s.Code() == 4 {
			return h.Transfer(ctx, conn, clientReq)
			//}
		}
		return nil, err
	}

	glog.Debugf("Transfer: resp[%s]", prototext.Format(resp))

	// TODO：这里获取server的方式可以改变下，放到Connect获取
	if conn.Server == "" {
		// md, ok := manager.RPCGameClient.GetCtxMetadata(ctx)
		// if ok {
		// 	conn.Server = md.Server()
		// }
	}

	return resp.GameResp, nil
}

func (h *Handler) Close(ctx context.Context, conn *base.Conn, clientReq *pb.GameReq) (*pb.GameResp, error) {
	conn.RLock()
	defer conn.RUnlock()

	req := &pb.CloseReq{
		// Uid:     conn.GuildID,
		// GameReq: clientReq,
	}

	// balancerCtx := module.NewBalancerCtx(
	// 	conn.UID,
	// 	conn.Server,
	// )

	_, err := manager.RPCGameClient.Close(ctx, req)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (h *Handler) Error(conn *base.Conn, err error, clientReq *pb.GameReq) (*pb.GameResp, error) {
	glog.Debugf("error: %v", err)
	errResp := &pb.S2CError{
		Error: int32(errors.Code(err)),
		Arg:   "",
	}

	errProto, err := proto.Marshal(errResp)
	if err != nil {
		glog.Errorf("proto marshal Error Error: %v", err)
		return nil, err
	}

	return &pb.GameResp{
		Command:   23,
		Serial:    clientReq.Serial,
		Timestamp: time.Now().Unix(),
		Content:   errProto,
	}, nil
}
