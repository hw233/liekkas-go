package base

import (
	"testing"
)

func TestConfigManager(t *testing.T) {
	cm := &ConfigManager{}
	cm.Init()
	cm.LoadConfig("../csv_files")

	t.Logf("config: %v", cm.CfgActionUnlockConfig.GetAllData())
}
