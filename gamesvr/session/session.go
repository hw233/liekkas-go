package session

import (
	"context"
	"database/sql"
	"sync"
	"time"

	"shared/protobuf/pb"
	"shared/utility/errors"
	"shared/utility/event"
	"shared/utility/glog"
	"shared/utility/param"
	"shared/utility/servertime"

	"shared/common"
	"shared/utility/safe"
	"shared/utility/session"

	"gamesvr/manager"
	"gamesvr/model"
)

const (
	// session status
	StatusInit    = 0 // 初始化
	StatusOnline  = 1 // 在线
	StatusOffline = 2 // 离线
)

type Builder struct{}

func (b *Builder) NewSession() session.Session {
	sess := &Session{
		EmbedManagedSession: &session.EmbedManagedSession{},
		User:                &model.User{},
	}

	return sess
}

type Session struct {
	sync.RWMutex
	*session.EmbedManagedSession
	*model.User

	portalSvr string
	// GuildServer  string
	// GuildContext *balancer.Context

	Token  string
	CTime  int64
	Status int8
	Serial int
}

func (s *Session) OnCreated(ctx context.Context, opts session.OnCreatedOpts) error {
	s.User = model.NewUser(opts.ID)

	err := manager.MySQL.Load(ctx, s.User)
	if err != nil {
		if err == sql.ErrNoRows {
			if opts.AllowNil {
				// 建号时初始设置
				err = s.User.InitForCreate(ctx)
				if err != nil {
					glog.Errorf(" user  InitForCreate error: %v", err)
					return err
				}

				err = manager.MySQL.Create(ctx, s.User)
				if err != nil {
					glog.Errorf("new user error: %v", err)
					return err
				}
			} else {
				return common.ErrUserNotFound
			}
		} else {
			glog.Errorf("load user error: %v", err)
			return err
		}
	}

	glog.Debugf("load user: %+v", *s.User)

	s.User.Init(ctx)

	s.ScheduleCall(5*time.Minute, func() {
		s.UpdateUserPower()

		// 玩家数据落地
		safe.Exec(5, func(i int) error {
			s.Lock()
			defer s.Unlock()

			err := manager.MySQL.Save(context.Background(), s.User)
			if err != nil {
				glog.Errorf("save user error: %v, times: %d ,\n", err, i)
			}

			return err
		})

		// 更新在线状态
		err := manager.Global.UserOnline(context.Background(), s.ID, servertime.Now().Unix(), 6*time.Minute)
		if err != nil {
			glog.Errorf("user online error: %v", err)
		}

		// 更新玩家每个sr和ssr的等级星级与战力到redis
		err = s.UserCharacterSimpleUpdate(context.Background())
		if err != nil {
			glog.Errorf("UserCharacterSimpleUpdate err: %v", err)
		}
	})

	s.ScheduleCall(30*time.Second, func() {
		s.Lock()
		defer s.Unlock()
		manager.EventQueue.ExecuteEventsInQueue(context.Background(), s.User.ID)
		// 每半分钟检查一下是否借到了角色
		s.MercenaryCharacterList(context.Background())

		err := s.GuildDataRefresh(context.Background())
		if err != nil {

		}

		err = s.PushGreetings(context.Background())
		if err != nil {
			glog.Errorf("push greetings error: %v\n", err)
		}
		err = s.PushGuildTasks(context.Background())
		if err != nil {
			glog.Errorf("push guild task items error: %v\n", err)
		}

		s.User.On30Second(servertime.Now().Unix())
	})

	s.ScheduleCall(time.Second, func() {
		s.Lock()
		defer s.Unlock()

		s.User.OnSecond(servertime.Now().Unix())
		s.TimerPush()
	})

	s.ScheduleCall(time.Hour, func() {
		s.Lock()
		defer s.Unlock()
		s.User.OnHour()
	})

	// 注册redis事件
	s.RegisterEventQueue(ctx)

	// 处理上线逻辑
	s.Online(ctx)

	return nil
}

// RegisterEventQueue 注册redis事件并执行
func (s *Session) RegisterEventQueue(ctx context.Context) {
	//
	event.UserEventHandler.Register(s.User.ID, common.EventTypeGraveyardHelp, func(Param *param.Param) error {
		HelpType, err := Param.GetInt(0)
		if err != nil {
			return errors.WrapTrace(err)
		}
		BuildUid, err := Param.GetInt64(1)
		if err != nil {
			return errors.WrapTrace(err)
		}
		Sec, err := Param.GetInt32(2)
		if err != nil {
			return errors.WrapTrace(err)
		}
		HelpAt, err := Param.GetInt64(3)
		if err != nil {
			return errors.WrapTrace(err)
		}
		s.User.GraveyardReceiveHelp(HelpType, BuildUid, Sec, HelpAt)
		return nil
	})

	event.UserEventHandler.Register(s.User.ID, common.EventTypeYggdrasilMail, func(Param *param.Param) error {
		fromUserName, err := Param.GetString(0)
		if err != nil {
			return errors.WrapTrace(err)
		}
		goodsId, err := Param.GetInt64(1)
		if err != nil {
			return errors.WrapTrace(err)
		}
		return s.User.Yggdrasil.TryTakeBackGoods(context.Background(), s.User.ID, fromUserName, goodsId)
	})

	event.UserEventHandler.Register(s.User.ID, common.EventTypeYggdrasilIntimacyChange, func(Param *param.Param) error {
		userId, err := Param.GetInt64(0)
		if err != nil {
			return errors.WrapTrace(err)
		}
		intimacy, err := Param.GetInt32(1)
		if err != nil {
			return errors.WrapTrace(err)
		}
		totalIntimacy, err := Param.GetInt32(2)
		if err != nil {
			return errors.WrapTrace(err)
		}

		s.User.AddYggPush(&pb.S2CYggdrasilIntimacyChange{
			UserId:        userId,
			IntimacyValue: intimacy,
			TotalIntimacy: totalIntimacy,
		})
		return nil
	})

	// 执行
	manager.EventQueue.ExecuteEventsInQueue(ctx, s.User.ID)

}

func (s *Session) OnTriggered(ctx context.Context) {

}

func (s *Session) OnClosed() {
	s.Lock()
	defer s.Unlock()

	glog.Debugf("OnClosed(): user: %d", s.ID)
	// 处理下线逻辑
	s.Offline()

	ctx := context.Background()

	safe.Exec(5, func(i int) error {
		err := manager.MySQL.Save(ctx, s.User)
		if err != nil {
			glog.Errorf("save user error: %v, times: %d ,\n", err, i)
		}

		return err
	})

	safe.Exec(5, func(i int) error {
		return manager.RPCServer.DelRecord(ctx, s.ID)
	})

	safe.Exec(5, func(i int) error {
		return manager.RPCServer.DecrBalance(ctx)
	})
}

func (s *Session) PrepareContext(ctx context.Context) context.Context {
	// TODO： 记录server减少redis压力
	// ctx = manager.RPCGuildClient.WithSinglecastCtx(ctx, grpc.SinglecastOpt(s.Guild.GuildID, ""))
	// ctx = manager.RPCPortalClient.WithSinglecastCtx(ctx, grpc.SinglecastOpt(s.ID, s.portalSvr))

	return ctx
}

func (s *Session) Defer(ctx context.Context) {
	s.MessagePush()
}

func (s *Session) RefreshPortalServer(server string) {
	s.portalSvr = server
}

// session存在即为上线
func (s *Session) Online(ctx context.Context) {
	// s.ConnToken = token
	// s.PortalServer = name

	now := servertime.Now().Unix()

	s.Status = StatusOnline

	s.User.OnOnline(now)
	// return nil

	// 更新在线状态
	err := manager.Global.UserOnline(ctx, s.ID, now, 6*time.Minute)
	if err != nil {
		glog.Errorf("UserOnline() error: %v", err)
	}

	// 更新公会玩家登录状态
	err = s.GuildSync(ctx, common.UserOnline)
	if err != nil {
		glog.Errorf("GuildSync() error: %v", err)
	}
}

// session关闭即为下线
func (s *Session) Offline() {
	// s.ConnToken = ""
	// s.PortalServer = ""
	s.Status = StatusOffline

	s.User.OnOffline()

	// 更新在线状态
	err := manager.Global.UserOffline(context.Background(), s.ID)
	if err != nil {
		glog.Errorf("user online error: %v", err)
	}

	// 更新公会玩家登录状态
	err = s.GuildSync(
		// make ctx
		context.Background(),
		// manager.RPCGuildClient.WithSinglecastCtx(
		// 	context.Background(),
		// 	grpc.SinglecastOpt(s.Guild.GuildID, ""),
		// ),
		common.UserOffline,
	)
	if err != nil {
		glog.Errorf("GuildSync() error: %v", err)
	}
}

// func (s *Session) Close() {
// 	// save to db
// 	err := manager.MySQL.Save(s.User)
// 	if err != nil {
// 		log.Printf("save user,%v\n", err)
// 	}
// }
//
// func (s *Session) Online(token, name string) {
// 	s.ConnToken = token
// 	s.ServerName = name
// 	s.Status = StatusOnline
// }
//
// func (s *Session) Offline() {
// 	s.ConnToken = ""
// 	s.ServerName = ""
// 	s.Status = StatusOffline
// }

func (s *Session) CheckToken(token string, ctime int64) error {
	if s.Token == "" {
		// replace new token
		s.Token = token
		s.CTime = ctime
		return nil
	}

	if s.Token != token {
		if ctime <= s.CTime {
			return common.ErrUserLoginInOtherClient
		}

		// replace new token
		s.Token = token
		s.CTime = ctime
	}

	return nil
}
