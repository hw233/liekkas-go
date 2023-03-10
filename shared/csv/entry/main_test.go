package entry

import (
	"log"
	"path/filepath"
	"shared/utility/glog"
	"testing"

	"shared/csv/base"
	"shared/utility/errors"
)

var CSV *Manager

func TestMain(m *testing.M) {
	glog.InitLog(&glog.Options{
		LogLevel:      "DEBUG",
		LogDir:        "./log",
		FileSizeLimit: 4294967296,
	})
	csvBase := &base.ConfigManager{}
	csvBase.Init()
	csvBase.LoadConfig(CSVPath)
	CSV = NewManager()
	err := CSV.Reload(csvBase)
	if err != nil {
		log.Printf("%+v", errors.Format(err))
	}

	m.Run()
}

var (
	CSVPath = filepath.Join("..", "data")
)
