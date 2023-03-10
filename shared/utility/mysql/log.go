package mysql

import "log"

type DBLogger interface {
	Debug(args ...interface{})
	Debugf(template string, args ...interface{})

	Info(args ...interface{})
	Infof(template string, args ...interface{})

	Error(args ...interface{})
	Errorf(template string, args ...interface{})
}

var dbLog DBLogger = &StdLogger{}

func SetLogger(logger DBLogger) {
	dbLog = logger
}

type StdLogger struct {
	DBLogger
}

func (std *StdLogger) Debug(args ...interface{}) {
	log.Print(args...)
}
func (std *StdLogger) Debugf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (std *StdLogger) Info(args ...interface{}) {
	log.Print(args...)
}
func (std *StdLogger) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func (std *StdLogger) Error(args ...interface{}) {
	log.Print(args...)
}
func (std *StdLogger) Errorf(template string, args ...interface{}) {
	log.Printf(template, args...)
}
