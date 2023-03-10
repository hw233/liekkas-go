package dblog

import (
	"shared/utility/glog"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type GDBLogger struct {
	logger *zap.Logger
}

func NewDBLogger() *GDBLogger {
	dbWriter := zapcore.Lock(glog.CreateDefaultLogFileStream("db"))
	dbErrorWritter := zapcore.Lock(glog.CreateDefaultLogFileStream("db_error"))

	dbLogger := glog.NewLogger([]zapcore.WriteSyncer{dbWriter},
		[]zapcore.WriteSyncer{dbErrorWritter, dbWriter, glog.Stderr}, 1)

	return &GDBLogger{
		logger: dbLogger,
	}
}

func (gl *GDBLogger) Debug(args ...interface{}) {
	gl.logger.Sugar().Debug(args...)
}
func (gl *GDBLogger) Debugf(template string, args ...interface{}) {
	gl.logger.Sugar().Debugf(template, args...)
}

func (gl *GDBLogger) Info(args ...interface{}) {
	gl.logger.Sugar().Info(args...)
}
func (gl *GDBLogger) Infof(template string, args ...interface{}) {
	gl.logger.Sugar().Infof(template, args...)
}

func (gl *GDBLogger) Error(args ...interface{}) {
	gl.logger.Sugar().Error(args...)
}
func (gl *GDBLogger) Errorf(template string, args ...interface{}) {
	gl.logger.Sugar().Errorf(template, args...)
}
