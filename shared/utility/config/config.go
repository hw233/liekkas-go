package config

import (
	"flag"
	"log"

	"github.com/spf13/viper"
)

var (
	configPath = flag.String("config-path", "./config.toml", "path of config file")
)

func SetConfigPath(path string) {
	*configPath = path
}

func Reload(config interface{}) error {
	err := load(config)
	if err != nil {
		return err
	}

	return nil
}

func Load(config interface{}) error {
	if !flag.Parsed() {
		flag.Parse()
	}

	// load config
	viper.SetConfigFile(*configPath)

	err := load(config)
	if err != nil {
		return err
	}

	return nil
}

func load(config interface{}) error {
	err := viper.ReadInConfig()
	if err != nil {
		log.Printf("read config error: %v", err)
		return err
	}

	err = viper.Unmarshal(config)
	if err != nil {
		log.Printf("unmarshal config error: %v", err)
		return err
	}

	return nil
}
