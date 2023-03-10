package controller

import (
	"context"
	"errors"
	"gamesvr/manager"
	"gamesvr/session"
	"log"
	"shared/common"
	"shared/protobuf/pb"
	customErr "shared/utility/errors"
	"shared/utility/glog"
	"shared/utility/router"
	"shared/utility/safe"
	"shared/utility/whitelist"
	"time"

	"google.golang.org/protobuf/proto"
)

type GameHandler struct {
	*pb.UnimplementedGameServer
	*router.Router
	WLSwitch       bool // 白名单开关
	MaintainSwitch bool // 维护开关
	*whitelist.MultiWhiteList
}

func NewGameHandler() (*GameHandler, error) {
	r, err := NewGameRoute()
	if err != nil {
		glog.Errorf("init router error: %v", err)
		return nil, err
	}

	MultiWhiteList := whitelist.NewMultiWhiteList()
	ctx := context.Background()
	err = reloadWhiteList(ctx, MultiWhiteList)
	if err != nil {
		return nil, err
	}
	wlSwitch, err := manager.Global.FetchWLSwitch(ctx)
	if err != nil {
		return nil, err
	}
	maintainSwitch, err := manager.Global.FetchMaintainSwitch(ctx)
	if err != nil {
		return nil, err
	}
	return &GameHandler{
		Router:         r,
		WLSwitch:       wlSwitch,
		MultiWhiteList: MultiWhiteList,
		MaintainSwitch: maintainSwitch,
	}, nil
}

func (gh *GameHandler) Connect(ctx context.Context, req *pb.ConnectReq) (*pb.ConnectResp, error) {
	defer safe.Recover()
	log.Printf("connect----------------> req: %+v", req)
	if req == nil || req.GameReq == nil {
		glog.Warnf("ERROR: nil req")
		return nil, errors.New("nil req")
	}

	if req.GameReq.Command == 1017 {
		resp, err := gh.ReconnectHandle(ctx, req)
		if err != nil {
			return newConnectResp(0, newErrResp(req.GameReq, err)), nil
		}
		return resp, nil

	} else if req.GameReq.Command == 1001 {
		resp, err := gh.LoginHandle(ctx, req)
		if err != nil {
			return newConnectResp(0, newErrResp(req.GameReq, err)), nil

		}
		return resp, nil
	}

	glog.Warnf("ERROR: illegal Command %d", req.GameReq.Command)
	return nil, common.ErrConnectIllegalCmd
}

func (gh *GameHandler) Transfer(ctx context.Context, req *pb.TransferReq) (*pb.TransferResp, error) {
	defer safe.Recover()
	log.Printf("transfer---------------->req: %+v", req)

	// check nil
	if req == nil || req.GameReq == nil {
		return nil, errors.New("nil req")
	}

	err := gh.PreHandle(req.Uid)
	if err != nil {
		glog.Warnf("gh.PreHandle error: %v", err)
		return newTransferResp(newErrResp(req.GameReq, err)), nil
	}

	out, err := gh.Call(ctx, req.Uid, req.GameReq.Command, req.GameReq.Content)
	if err != nil {
		glog.Warnf("gh.Call error: %v", err)
		return newTransferResp(newErrResp(req.GameReq, err)), nil
	}

	return newTransferResp(newGameResp(req.GameReq, out)), nil
}

func (gh *GameHandler) Call(ctx context.Context, uid int64, command int32, in []byte) ([]byte, error) {
	f, ok := gh.Route(command)
	if !ok {
		glog.Warnf("GameHandler.Call: not found in router,cmd:%d", command)
		return nil, errors.New("not found in router")
	}

	// get protobuf request
	in2 := f.In(2).(proto.Message)
	err := proto.Unmarshal(in, in2)
	if err != nil {
		glog.Warnf("GameHandler.Call: proto.Unmarshal error: %v", err)
		return nil, err
	}

	// get session
	managedSess, err := manager.SessManager.GetSession(ctx, uid)
	if err != nil {
		return nil, err
	}

	sess, ok := managedSess.(*session.Session)
	if !ok {
		return nil, errors.New("invalid session")
	}

	glog.Debugf("player [%d] request %s: \n%v\n", uid, in2.ProtoReflect().Descriptor().Name(), in2)

	sess.Lock()
	defer sess.Unlock()

	ctx = sess.PrepareContext(ctx)

	// call
	ret := f.Call(sess, ctx, in2)

	sess.Defer(ctx)

	// handle error
	err, ok = ret[1].(error)
	if ok && err != nil {
		glog.Warnf("GameHandler.Call: return error: %+v", customErr.Format(err))
		return nil, err
	}

	// handle response
	out, err := proto.Marshal(ret[0].(proto.Message))
	if err != nil {
		glog.Warnf("GameHandler.Call: Marshal error: %v", err)
		return nil, err
	}

	glog.Debugf("player [%d] respone %s: \n%v\n", uid, ret[0].(proto.Message).ProtoReflect().Descriptor().Name(), ret[0])

	return out, nil
}

// func (h *Handler) SessionCall(ctx context.Context, sess *session.CommonSession, f interface{}) error {
// 	// get type of function
// 	t := reflect.TypeOf(f)
//
// 	// get type of the second of function parameter
// 	ti := t.In(1).Elem()
//
// 	// new request struct and unmarshal from protobuf bytes
// 	i := reflect.New(ti).Interface()
// 	err := gh.UnmarshalProto(i.(proto.ServiceMessage))
// 	if err != nil {
// 		return nil
// 	}
//
// 	// call function
// 	val := reflect.ValueOf(f).Call([]reflect.Operation{
// 		reflect.ValueOf(sess),
// 		reflect.ValueOf(i),
// 	})
//
// 	// handle return value
// 	// val0 = response, val1 = error
// 	val0, val1 := val[0], val[1]
//
// 	// handle error
// 	err, ok := val1.Interface().(error)
// 	if ok && err != nil {
// 		return err
// 	}
//
// 	// handle response
// 	m, ok := val0.Interface().(proto.ServiceMessage)
// 	if !ok {
// 		return fmt.Errorf("type of resp not proto.ServiceMessage")
// 	}
//
// 	err = gh.MarshalProto(m)
// 	if err != nil {
// 		return err
// 	}
//
// 	return nil
// }

func (gh *GameHandler) NewGroupMailNotify(ctx context.Context, req *pb.NewGroupMailReq) (*pb.NewGroupMailResp, error) {
	voMail := req.Mail
	sendUsers := req.Mail.Users
	if voMail.SendAll {
		sendUsers = manager.SessManager.GetAllSessionId()
	}

	mail, err := common.ParseFromVOServerMail(voMail.Mail)
	if err != nil {
		return nil, err
	}

	for _, userId := range sendUsers {
		if !manager.SessManager.HasSession(userId) {
			continue
		}

		managedSession, err := manager.SessManager.GetSession(ctx, userId)
		if err != nil {
			glog.Errorf("send group mail [%d] to user [%d] error: %+v", mail.Id, userId, err)
			continue
		}

		session, ok := managedSession.(*session.Session)
		if !ok {
			glog.Errorf("send group mail [%d] to user [%d] error: invalid session", mail.Id, userId)
			continue
		}

		session.Lock()
		session.User.AddServerGroupMail(mail)
		session.Unlock()
	}

	return &pb.NewGroupMailResp{}, nil
}

func (gh *GameHandler) NewPersonalMailNotify(ctx context.Context, req *pb.NewPersonalMailReq) (*pb.NewPersonalMailResp, error) {
	mail, err := common.ParseFromVOServerMail(req.Mail.Mail)
	if err != nil {
		return nil, err
	}

	userId := req.Mail.User

	managedSession, err := manager.SessManager.GetSession(ctx, userId)
	if err != nil {
		glog.Errorf("send personal mail [%d] to user [%d] error: %+v", mail.Id, userId, err)
		return nil, err
	}

	session, ok := managedSession.(*session.Session)
	if !ok {
		glog.Errorf("send group mail [%d] to user [%d] error: invalid session", mail.Id, userId)
		return nil, errors.New("invalid session")
	}

	session.Lock()
	session.User.AddServerPersonalMail(mail)
	session.Unlock()

	return &pb.NewPersonalMailResp{}, nil
}

func (gh *GameHandler) ReloadWhiteList(ctx context.Context, req *pb.ReloadWhiteListReq) (*pb.ReloadWhiteListResp, error) {
	// 刷新一下
	err := reloadWhiteList(ctx, gh.MultiWhiteList)
	if err != nil {
		return nil, err
	}
	wlSwitch, err := manager.Global.FetchWLSwitch(ctx)
	if err != nil {
		return nil, err
	}
	gh.WLSwitch = wlSwitch
	return &pb.ReloadWhiteListResp{}, nil

}

func reloadWhiteList(ctx context.Context, multi *whitelist.MultiWhiteList) error {
	defer safe.Recover()
	// 刷新一下
	ops, err := manager.Global.GetIdWhiteList(ctx)
	if err != nil {
		return err
	}

	multi.Reload(ops)
	return nil
}

func (gh *GameHandler) ReloadMaintain(ctx context.Context, req *pb.ReloadMaintainReq) (*pb.ReloadMaintainResp, error) {
	gh.reloadMaintain(ctx)
	return &pb.ReloadMaintainResp{}, nil

}

func (gh *GameHandler) reloadMaintain(ctx context.Context) {
	maintainSwitch, err := manager.Global.FetchMaintainSwitch(ctx)
	if err != nil {
		glog.Errorf("FetchMaintainSwitch FetchMaintainSwitch err:%v", err)
		maintainSwitch = manager.CSV.GlobalEntry.SpareMaintainSwitch
	}
	if gh.MaintainSwitch == maintainSwitch {
		return
	}
	if maintainSwitch {
		manager.SessManager.Close()
	}

	gh.MaintainSwitch = maintainSwitch
}

func (gh *GameHandler) PreHandle(uid int64) error {
	if gh.MaintainSwitch && gh.WLSwitch {
		if !gh.Filter(whitelist.NewOpsWithUid(uid)) {
			glog.Debugf("PreHandle Intercepted ,uid:%d", uid)
			return common.ErrSeverMaintaining
		}

	}

	return nil
}

func (gh *GameHandler) ReloadAnnouncement(ctx context.Context, req *pb.ReloadAnnouncementReq) (*pb.ReloadAnnouncementResp, error) {
	err := manager.Announcements.LoadAnnouncements(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ReloadAnnouncementResp{}, nil
}

func (gh *GameHandler) Init() {
	manager.Timer.ScheduleFunc(time.Minute, func() {
		gh.reloadMaintain(context.Background())
	})

}
