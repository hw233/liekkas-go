package main

import (
	"log"
	"os"
)

type managerGenerator struct {
	configs []*generator
}

func NewConfigsManagerMaker(configs []*generator) *managerGenerator {
	return &managerGenerator{
		configs: configs,
	}
}

// fileOut is .go file output directory.
func (m *managerGenerator) Generate(fileOut string, fileCfg string) {
	fileName := fileOut + "/" + "config_manager" + ".go"
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0644)
	if err != nil {
		log.Printf("open file %s error: %v", fileName, err)
	}

	defer file.Close()

	var con string
	con += "package base"
	con += "\n\n"
	// con += "import \"log\""
	con += "import \"shared/utility/glog\""
	con += "\n\n"

	// declare ConfigManager struct.
	con += "type ConfigManager struct {"
	con += "\n"

	maxNameLen := 0
	for i := 0; i < len(m.configs); i++ {
		nameLen := len(m.configs[i].typeName + "Config")
		if maxNameLen < nameLen {
			maxNameLen = nameLen
		}
	}

	for i := 0; i < len(m.configs); i++ {
		con += addSpace("", 4)
		name := m.configs[i].typeName + "Config"
		con += name + addSpace("", maxNameLen-len(name)+1) + "*" + name
		con += "\n"
	}
	con += "}"
	con += "\n\n"

	// Init func
	con += "func (m *ConfigManager) Init() {"
	con += "\n"
	for i := 0; i < len(m.configs); i++ {
		con += addSpace("", 4)
		name := m.configs[i].typeName + "Config"
		con += "m." + name + " = " + "New" + name + "()"
		con += "\n"
	}
	con += "}"
	con += "\n\n"

	// add set/reload function.
	addReloadFunc := func(typeName, fileName string) {
		// SetConfig
		con += "// set " + typeName + " interface \n"
		con += "func (m *ConfigManager) "
		con += "Set" + typeName + "Config(config *"
		con += typeName + "Config) {"
		con += "\n"

		con += addSpace("", 4)
		con += "m." + typeName + "Config = config"
		con += "\n"

		con += "}"
		con += "\n\n"

		// ReloadConfig
		con += "// reload config " + typeName + " interface \n"
		con += "func (m *ConfigManager) "
		con += "Reload" + typeName + "Config() *" + typeName + "Config {"
		con += "\n"
		con += addSpace("", 4)
		con += "reloadConfig := New" + typeName + "Config()"
		con += "\n"

		con += addSpace("", 4)
		con += "if !reloadConfig.Load(\"" + fileCfg + "/" + lower(fileName) + ".csv\") {"
		con += "\n"
		con += addSpace("", 8)
		con += "glog.Info(" + "\"Load " + fileCfg + "/" + lower(fileName) + ".csv" + " error\")"
		con += "\n"
		con += addSpace("", 8)
		con += "return nil"
		con += "\n"
		con += addSpace("", 4)
		con += "}"
		con += "\n"

		con += addSpace("", 4)
		con += "return reloadConfig"
		con += "\n"

		con += "}"
		con += "\n\n"
	}

	// Set detail config function
	for i := 0; i < len(m.configs); i++ {
		addReloadFunc(m.configs[i].typeName, m.configs[i].fileName)
	}
	con += "\n"

	// LoadConfig func
	con += "func (m *ConfigManager) LoadConfig(path string) bool {"
	con += "\n"

	con += addSpace("", 4)
	con +=
		`config2Loader := [] struct {
		path   string
		loader Config
	} {`
	con += "\n"

	// add all file path, and its config variable.
	for i := 0; i < len(m.configs); i++ {
		con += addSpace("", 8)
		con += "{ \"" + lower(m.configs[i].fileName) + ".csv\"," + "m." + m.configs[i].typeName + "Config" + " },"
		con += "\n"
	}
	con += addSpace("", 4)
	con += `}
	// range config map.
	for i := 0; i < len(config2Loader); i++ {
		if !config2Loader[i].loader.Load(path + "/" + config2Loader[i].path) {
			glog.Info("Load ", config2Loader[i].path, " file error")
			return false
		} else {
			glog.Info("Load ", config2Loader[i].path, " file success")
		}
	}
	return true`

	con += "\n"
	con += "}"
	con += "\n\n"

	// global variable
	con += "var ConfigManagerObj ConfigManager"
	con += "\n\n"
	con += `func init() {
	ConfigManagerObj.Init()
}`
	con += "\n"

	// write to file
	_, err = file.WriteString(con)
	if err != nil {
		log.Printf("write string to %s error: %v", fileName, err)
	}
}
