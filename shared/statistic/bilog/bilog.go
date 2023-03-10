package bilog

import (
	"path"
	"shared/utility/glog"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var eventLogger *zap.Logger
var snapshotLogger *zap.Logger

type LogObjMarshaler zapcore.ObjectMarshaler

type EventData struct {
	key  string
	data LogObjMarshaler
}

func NewEventData(key string, data LogObjMarshaler) *EventData {
	return &EventData{
		key:  key,
		data: data,
	}
}

func (ed *EventData) MarshalLogObject(encoder zapcore.ObjectEncoder) error {
	return ed.data.MarshalLogObject(encoder)
}

func log(logger *zap.Logger, commonDatas []LogObjMarshaler, eventDatas []*EventData) {
	if logger == nil {
		return
	}

	if commonDatas == nil {
		commonDatas = []LogObjMarshaler{}
	}

	if eventDatas == nil {
		eventDatas = []*EventData{}
	}

	eventDataFields := make([]zap.Field, 0, len(commonDatas)+len(eventDatas))
	for _, commonData := range commonDatas {
		eventDataFields = append(eventDataFields, zap.Inline(commonData))
	}

	for _, eventData := range eventDatas {
		eventDataFields = append(eventDataFields, zap.Object(eventData.key, eventData))
	}

	logger.Info("", eventDataFields...)
}

func InitEventLogger(appName, dir string) {
	spliter := NewBILogSpliter()

	prefix := strings.Join([]string{"bilievent", appName}, "_")
	fileStream := glog.NewLogFileStream(path.Join(dir, "event"), prefix, spliter)

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{})
	writer := zapcore.Lock(fileStream)
	levelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, levelEnablerFunc),
	)

	eventLogger = zap.New(core)
}

func EventLog(commonDatas []LogObjMarshaler, eventDatas []*EventData) {
	log(eventLogger, commonDatas, eventDatas)
}

func InitSnapshotLogger(appName, dir string) {
	spliter := NewBILogDailySpliter()
	prefix := strings.Join([]string{"biliuser", appName}, "_")
	fileStream := glog.NewLogFileStream(path.Join(dir, "user"), prefix, spliter)

	encoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{})
	writer := zapcore.Lock(fileStream)
	levelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return true
	})

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, writer, levelEnablerFunc),
	)

	snapshotLogger = zap.New(core)

}

func SnapshotLog(commonDatas []LogObjMarshaler, eventDatas []*EventData) {
	log(snapshotLogger, commonDatas, eventDatas)
}
