package session

import (
	"context"
	"time"

	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"

	"gamesvr/manager"
	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/servertime"
)

func (s *Session) pushTest() error {
	for i := 0; i < 10; i++ {
		time.Sleep(5 * time.Second)

		resp := &pb.TestResp{
			Uid:  s.ID,
			Name: s.portalSvr,
		}

		err := s.push(1, resp)
		if err != nil {
			return errors.WrapTrace(err)
		}
	}

	return nil
}

func (s *Session) pushSomething() error {
	err := s.push(1, nil)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return nil
}

func (s *Session) MessagePush() {
	if s.User == nil {
		return
	}
	s.TryPushQuestUpdateNotify()
	s.TryPushExplore()
	s.TryPushTowerUpdateNotify()
	s.TryPushMailNotify()
	s.TryPushYgg()
	s.TryPushStoreNotify()
}

func (s *Session) TimerPush() {
	if s.User == nil {
		return
	}

	s.TryPushLevelNotify()
	s.TryPushQuestUpdateNotify()
	s.DailyRefreshNotify()
	s.TryPushScorePassNotify()
	s.TryPushYgg()
	s.TryPushGraveyard()
}

// 推送
func (s *Session) push(command int32, message proto.Message) error {
	pushReq, err := newPushReq(s.ID, command, message)
	if err != nil {
		return errors.WrapTrace(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// ctx = manager.RPCPortalClient.WithSinglecastCtx(ctx, grpc.SinglecastOpt(s.ID, s.portalSvr))

	glog.Debugf("player [%d] push %s: \n%v\n", s.User.ID, message.ProtoReflect().Descriptor().Name(),
		prototext.Format(message))

	_, err = manager.RPCPortalClient.Push(ctx, pushReq)
	if err != nil {
		return errors.WrapTrace(err)
	}

	return err
}

func newPushReq(uid int64, commend int32, message proto.Message) (*pb.PushReq, error) {
	content, err := proto.Marshal(message)
	if err != nil {
		return nil, err
	}

	return &pb.PushReq{
		GameResp: &pb.GameResp{
			Command:   commend,
			Serial:    0,
			Timestamp: servertime.Now().Unix(),
			Content:   content,
		},
		Uid: uid,
	}, nil
}
