package manager

import (
	"shared/statistic/bilog"
	"shared/utility/servertime"
	"time"
)

func BIInit() {
	bilog.InitEventLogger(Conf.AppName, Conf.BILogDir)
	bilog.InitSnapshotLogger(Conf.AppName, Conf.BILogDir)

	Timer.ScheduleFunc(time.Minute, BIOnlineCount)
}

func BIOnlineCount() {
	now := servertime.Now()

	logId := bilog.FormatGameLogId(Conf.AppName, bilog.EventNameOnlineCount, now.Unix())
	logData := &bilog.GameOnlineUsers{
		Time:        servertime.Now(),
		LogId:       logId,
		GameBaseId:  "",
		GameId:      0,
		Platform:    "",
		ZoneId:      0,
		ChannelId:   0,
		OnlineCount: SessManager.GetSessionCount(),
	}

	bilog.SnapshotLog([]bilog.LogObjMarshaler{logData}, nil)
}
