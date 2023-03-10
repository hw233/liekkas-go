package entry

import (
	"reflect"

	"foreplay/csv/base"
	"shared/utility/errors"
	"shared/utility/glog"
)

type Manager struct {
	Patches *Patches
	AppList *AppList
}

func NewManager() *Manager {
	return &Manager{
		Patches: NewPatches(),
		AppList: NewAppList(),
	}
}

type Config struct {
	*base.ConfigManager
}

type Entry interface {
	Reload(config *Config) error
}

func (m *Manager) Reload(config *base.ConfigManager) error {
	v := reflect.ValueOf(m).Elem()

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		data := field.Interface()

		entry, ok := data.(Entry)
		if ok {
			err := entry.Reload(&Config{
				ConfigManager: config,
			})
			if err != nil {
				glog.Errorf("ERROR: load %s error: %+v", field.Type(), errors.Format(err))
				return errors.WrapTrace(err)
			}
		}
	}

	return nil
}
