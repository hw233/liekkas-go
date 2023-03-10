package glog

import (
	"flag"
	"os"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logLevel      zapcore.Level = zapcore.InfoLevel
	logRootDir    string        = "./log"
	fileSizeLimit int64         = 0

	startupLogRootDir    *string = flag.String("log-dir", "", "log root dir, eg. --log-dir ./")
	startupFileSizeLimit *int64  = flag.Int64("log-size", -1, "single log file size limit, splited while size is over, eg. --log-size /appdir/log")

	Stdout            = zapcore.Lock(os.Stdout)
	Stderr            = zapcore.Lock(os.Stderr)
	RuntimeFileWriter zapcore.WriteSyncer
	ErrorFileWriter   zapcore.WriteSyncer

	defaultLogger *zap.Logger = NewLogger([]zapcore.WriteSyncer{Stdout}, []zapcore.WriteSyncer{Stderr}, 1)
	runtimeLogger *zap.Logger = defaultLogger
)

type Options struct {
	LogLevel      string
	LogDir        string
	FileSizeLimit int64
}

func InitLog(opts *Options) error {
	if !flag.Parsed() {
		flag.Parse()
	}

	err := LoadOptions(opts)
	if err != nil {
		return err
	}

	if *startupLogRootDir != "" {
		logRootDir = *startupLogRootDir
	}

	if *startupFileSizeLimit != -1 {
		fileSizeLimit = *startupFileSizeLimit
	}

	RuntimeFileWriter = zapcore.Lock(CreateDefaultLogFileStream("runtime"))
	ErrorFileWriter = zapcore.Lock(CreateDefaultLogFileStream("error"))

	runtimeLogger = NewLogger([]zapcore.WriteSyncer{RuntimeFileWriter, Stdout},
		[]zapcore.WriteSyncer{ErrorFileWriter, RuntimeFileWriter, Stderr}, 1)

	return nil
}

func CreateDefaultLogFileStream(subDir string) *LogFileStream {
	spliter := NewLogFileHourSpliter(time.Now, int(fileSizeLimit))
	fileStream := NewLogFileStream(path.Join(logRootDir, subDir), "", spliter)

	return fileStream
}

func LoadOptions(opts *Options) error {
	if opts == nil {
		return nil
	}

	if opts.LogLevel != "" {
		err := SetLogLevel(opts.LogLevel)
		if err != nil {
			return err
		}
	}

	if opts.LogDir != "" {
		SetRootDir(opts.LogDir)
	}

	if opts.FileSizeLimit >= 0 {
		SetFileSizeLimit(opts.FileSizeLimit)
	}

	return nil
}

func SetLogLevel(level string) error {
	return logLevel.Set(level)
}

func SetRootDir(dir string) {
	logRootDir = dir
}

func SetFileSizeLimit(limit int64) {
	fileSizeLimit = limit
}

// Debug
func Debug(args ...interface{}) {
	runtimeLogger.Sugar().Debug(args...)
}

func Debugf(template string, args ...interface{}) {
	runtimeLogger.Sugar().Debugf(template, args...)
}

// Info
func Info(args ...interface{}) {
	runtimeLogger.Sugar().Info(args...)
}

func Infof(template string, args ...interface{}) {
	runtimeLogger.Sugar().Infof(template, args...)
}

// Warn
func Warn(args ...interface{}) {
	runtimeLogger.Sugar().Warn(args...)
}

func Warnf(template string, args ...interface{}) {
	runtimeLogger.Sugar().Warnf(template, args...)
}

// Error
func Error(args ...interface{}) {
	runtimeLogger.Sugar().Error(args...)
}

func Errorf(template string, args ...interface{}) {
	runtimeLogger.Sugar().Errorf(template, args...)
}

// Panic
func Panic(args ...interface{}) {
	runtimeLogger.Sugar().Panic(args...)
}

func Panicf(template string, args ...interface{}) {
	runtimeLogger.Sugar().Panicf(template, args...)
}

// Fatal
func Fatal(args ...interface{}) {
	runtimeLogger.Sugar().Fatal(args...)
}

func Fatalf(template string, args ...interface{}) {
	runtimeLogger.Sugar().Fatalf(template, args...)
}
