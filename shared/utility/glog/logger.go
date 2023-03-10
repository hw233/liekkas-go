package glog

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewCommonEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:          "T",
		MessageKey:       "M",
		LevelKey:         "L",
		CallerKey:        "C",
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		EncodeCaller:     zapcore.ShortCallerEncoder,
		ConsoleSeparator: " ",
	}
}

func NewErrorEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:          "T",
		MessageKey:       "M",
		LevelKey:         "L",
		StacktraceKey:    "S",
		EncodeTime:       zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"),
		EncodeLevel:      zapcore.CapitalLevelEncoder,
		ConsoleSeparator: " ",
	}
}

func NewLogger(commonWriters []zapcore.WriteSyncer, errorWriters []zapcore.WriteSyncer, callerSkip int) *zap.Logger {
	commonEncoder := zapcore.NewConsoleEncoder(NewCommonEncoderConfig())
	commonLevelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel && lvl < zap.ErrorLevel
	})

	errorEncoder := zapcore.NewConsoleEncoder(NewErrorEncoderConfig())
	errorLevelEnablerFunc := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= logLevel && lvl >= zap.ErrorLevel
	})

	core := zapcore.NewTee(
		zapcore.NewCore(commonEncoder, zap.CombineWriteSyncers(commonWriters...), commonLevelEnablerFunc),
		zapcore.NewCore(errorEncoder, zap.CombineWriteSyncers(errorWriters...), errorLevelEnablerFunc),
	)

	logger := zap.New(core, zap.AddStacktrace(zap.ErrorLevel), zap.AddCaller(), zap.AddCallerSkip(callerSkip))

	return logger
}
