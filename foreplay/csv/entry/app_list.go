package entry

import (
	"shared/common"
	"shared/utility/errors"
	"shared/utility/transfer"
	"sync"
)

type AppCfg struct {
	Channel            string   `json:"channel"`
	AppVersion         string   `json:"app_version"`
	AppUrl             string   `json:"app_url"`
	CompatibleVersions []string `json:"compatible_versions"`
}

type AppList struct {
	sync.RWMutex

	Apps map[string]*AppCfg
}

func NewAppList() *AppList {
	return &AppList{}
}

func (al *AppList) Reload(config *Config) error {
	al.Lock()
	defer al.Unlock()

	apps := map[string]*AppCfg{}

	for _, appCSV := range config.AppListConfig.GetAllData() {
		appCfg := &AppCfg{}

		err := transfer.Transfer(appCSV, appCfg)
		if err != nil {
			return errors.WrapTrace(err)
		}

		_, ok := apps[appCfg.Channel]
		if ok {
			return errors.WrapTrace(errors.New("duplicate channel: %s", appCfg.Channel))
		}

		apps[appCfg.Channel] = appCfg
	}

	al.Apps = apps

	return nil
}

func (al *AppList) GetAppCfg(channel string) (*AppCfg, error) {
	appCfg, ok := al.Apps[channel]
	if !ok {
		return nil, common.ErrAppCfgNotFound
	}

	return appCfg, nil
}
