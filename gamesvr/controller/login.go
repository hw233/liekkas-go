package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"

	"gamesvr/manager"
	"gamesvr/session"
	"shared/common"
	"shared/protobuf/pb"
	"shared/utility/glog"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

func CheckUserLogin(ctx context.Context, uid int64, token string) error {
	if uid == 0 {
		return fmt.Errorf("ERROR: user not found: %d", uid)
	}

	resp, err := manager.RPCLoginClient.CheckToken(ctx, &pb.CheckTokenReq{
		UserID: uid,
		Token:  token,
	})
	if err != nil {
		return err
	}

	if !resp.OK {
		return common.ErrConnectTokenInvalid
	}

	return nil
}

// 处理登录
func (gh *GameHandler) LoginHandle(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectResp, error) {
	loginReq := &pb.C2SLogin{}
	err := proto.Unmarshal(req.GameReq.Content, loginReq)
	if err != nil {
		glog.Errorf("ERROR: proto unmarshal error: %v", err)
		return nil, err
	}
	err = gh.PreHandle(loginReq.UserId)
	if err != nil {
		return nil, err
	}

	glog.Debugf("player [%d] request %s: \n%v\n", loginReq.UserId,
		loginReq.ProtoReflect().Descriptor().Name(), prototext.Format(loginReq))

	// 去login服务器判断玩家信息合法性
	// err = CheckUserLogin(ctx, loginReq.UserId, loginReq.Token)
	// if err != nil {
	// 	glog.Warnf("ERROR: check user login error: %v", err)
	// 	return nil, err
	// }

	// uidS := strconv.FormatInt(loginReq.UserId, 10)

	// 判断session是否在其他game，如果在的话这边就不初始化session了，下次请求会分发到对应的game
	err = manager.RPCServer.SetRecord(ctx, loginReq.UserId)
	if err != nil {
		glog.Errorf("SetRecord(%d) error: %v", loginReq.UserId, err)
		if err == redis.Nil {
			return nil, common.ErrConnectUserInOtherGameSvr
		}

		return nil, common.ErrDefault
	}

	// 初始化session
	managedSess, err := manager.SessManager.NewSessionIfNotExist(ctx, loginReq.UserId)
	if err != nil {
		glog.Errorf("get session error: %v", err)
		return nil, err
	}

	sess, ok := managedSess.(*session.Session)
	if !ok {
		glog.Errorf("invalid session")
		return nil, errors.New("invalid session")
	}

	sess.Lock()
	defer sess.Unlock()
	// sess, err := session.NewSessionIfNotExist(loginReq.UserId)
	// if err != nil {
	// 	glog.Warnf("ERROR: new session error: %v", err)
	// 	return nil, err
	// }
	//

	err = sess.CheckToken(req.Token, req.Ctime)
	if err != nil {
		glog.Errorf("invalid token")
		return nil, err
	}

	sess.RefreshPortalServer(req.Server)
	// sess.Online(req.Token, req.Server)
	//  检查token
	// err = sess.CheckAndUpdateToken(loginReq.ConnToken)
	// if err != nil {
	// 	glog.Warnf("ERROR: check token error: %v", err)
	// 	return nil, err
	// }

	loginResp := &pb.S2CLogin{
		IsNewUser: false,
		Token:     loginReq.Token,
		UserId:    sess.GetUserId(),
	}

	out, err := proto.Marshal(loginResp)
	if err != nil {
		glog.Errorf("proto marshal error: %v", err)
		return nil, err
	}

	glog.Debugf("loginResp: %+v", loginResp)

	return newConnectResp(loginReq.UserId, newGameResp(req.GameReq, out)), nil
}

// 处理重连
func (gh *GameHandler) ReconnectHandle(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectResp, error) {
	reconnectReq := &pb.C2SUserReconnect{}
	err := proto.Unmarshal(req.GameReq.Content, reconnectReq)
	if err != nil {
		glog.Warnf("ERROR: proto unmarshal error: %v", err)
		return nil, err
	}
	err = gh.PreHandle(reconnectReq.UserId)
	if err != nil {
		return nil, err
	}
	glog.Debugf("player [%d] request %s: \n%v\n", reconnectReq.UserId,
		reconnectReq.ProtoReflect().Descriptor().Name(), prototext.Format(reconnectReq))

	// 去login服务器判断玩家信息合法性
	// err = CheckUserLogin(ctx, reconnectReq.UserId, reconnectReq.Token)
	// if err != nil {
	// 	glog.Warnf("ERROR: check user login error: %v", err)
	// 	return nil, err
	// }

	// 判断session是否在其他game
	err = manager.RPCServer.SetRecord(ctx, reconnectReq.UserId)
	if err != nil {
		glog.Errorf("SetRecord(%d) error: %v", reconnectReq.UserId, err)
		if err == redis.Nil {
			return nil, common.ErrConnectUserInOtherGameSvr
		}

		return nil, common.ErrDefault
	}

	// Session是否在本gamesvr
	managedSess, err := manager.SessManager.NewSessionIfNotExist(ctx, reconnectReq.UserId)
	if err != nil {
		return nil, err
	}

	sess, ok := managedSess.(*session.Session)
	if !ok {
		return nil, errors.New("invalid session")
	}

	sess.Lock()
	defer sess.Unlock()

	// sess, err := session.NewSessionIfNotExist(reconnectReq.UserId)
	// if err != nil {
	// 	glog.Warnf("ERROR: new session error: %v", err)
	// 	return nil, err
	// }
	//

	err = sess.CheckToken(req.Token, req.Ctime)
	if err != nil {
		glog.Errorf("invalid token")
		return nil, err
	}

	sess.RefreshPortalServer(req.Server)
	// sess.Online(req.Token, req.Server)

	reconnectResp := &pb.S2CUserReconnect{
		Token: reconnectReq.Token,
	}

	out, err := proto.Marshal(reconnectResp)
	if err != nil {
		glog.Warnf("ERROR: proto marshal error: %v", err)
		return nil, err
	}

	glog.Debugf("reconnectResp: %+v", reconnectResp)

	return newConnectResp(reconnectReq.UserId, newGameResp(req.GameReq, out)), nil
}
