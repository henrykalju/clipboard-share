package common

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BackendUrl string
}

var conf *viper.Viper

const (
	BACKEND_URL_KEY     = "BACKEND_URL"
	BACKEND_URL_DEFAULT = "localhost:8080"

	CONF_DIR  = "clipbaord-share"
	CONF_FILE = "config.json"
)

func InitConfig() {
	conf = viper.New()

	setDefaults(conf)

	conf.SetConfigName("config")
	conf.SetConfigType("json")
	confPath, err := os.UserConfigDir()
	if err != nil {
		panic(err)
	}
	conf.AddConfigPath(filepath.Join(confPath, CONF_DIR, CONF_FILE))

	err = conf.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			panic(err)
		}
	}
}

func setDefaults(c *viper.Viper) {
	c.SetDefault(BACKEND_URL_KEY, BACKEND_URL_DEFAULT)
}

func GetConf() Config {
	return Config{
		BackendUrl: conf.GetString(BACKEND_URL_KEY),
	}
}

func SetConf(c Config) error {
	conf.Set(BACKEND_URL_KEY, c.BackendUrl)
	return conf.WriteConfig()
}
