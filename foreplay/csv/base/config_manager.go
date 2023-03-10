package base

import "shared/utility/glog"

type ConfigManager struct {
	AppListConfig   *AppListConfig
	PatchListConfig *PatchListConfig
}

func (m *ConfigManager) Init() {
	m.AppListConfig = NewAppListConfig()
	m.PatchListConfig = NewPatchListConfig()
}

// set AppList interface
func (m *ConfigManager) SetAppListConfig(config *AppListConfig) {
	m.AppListConfig = config
}

// reload config AppList interface
func (m *ConfigManager) ReloadAppListConfig() *AppListConfig {
	reloadConfig := NewAppListConfig()
	if !reloadConfig.Load("../../foreplay/csv/data/app_list.csv") {
		glog.Info("Load ../../foreplay/csv/data/app_list.csv error")
		return nil
	}
	return reloadConfig
}

// set PatchList interface
func (m *ConfigManager) SetPatchListConfig(config *PatchListConfig) {
	m.PatchListConfig = config
}

// reload config PatchList interface
func (m *ConfigManager) ReloadPatchListConfig() *PatchListConfig {
	reloadConfig := NewPatchListConfig()
	if !reloadConfig.Load("../../foreplay/csv/data/patch_list.csv") {
		glog.Info("Load ../../foreplay/csv/data/patch_list.csv error")
		return nil
	}
	return reloadConfig
}

func (m *ConfigManager) LoadConfig(path string) bool {
	config2Loader := []struct {
		path   string
		loader Config
	}{
		{"app_list.csv", m.AppListConfig},
		{"patch_list.csv", m.PatchListConfig},
	}
	// range config map.
	for i := 0; i < len(config2Loader); i++ {
		if !config2Loader[i].loader.Load(path + "/" + config2Loader[i].path) {
			glog.Info("Load ", config2Loader[i].path, " file error")
			return false
		} else {
			glog.Info("Load ", config2Loader[i].path, " file success")
		}
	}
	return true
}

var ConfigManagerObj ConfigManager

func init() {
	ConfigManagerObj.Init()
}
