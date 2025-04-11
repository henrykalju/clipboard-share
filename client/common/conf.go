package common

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	BackendUrl string
	Username   string
	Password   string
}

var conf *viper.Viper

const (
	BACKEND_URL_KEY     = "BACKEND_URL"
	BACKEND_URL_DEFAULT = "localhost:8080"

	USERNAME_KEY = "USERNAME"
	PASSWORD_KEY = "PASSWORD"

	CONF_PATH = "clipboard-share"
	CONF_FILE = "config.json"
)

func InitConfig() error {
	conf = viper.New()

	setDefaults(conf)

	conf.SetConfigName("config")
	conf.SetConfigType("json")
	confPath, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("error getting user config dir: %w", err)
	}
	confPath = filepath.Join(confPath, CONF_PATH)
	confFile := filepath.Join(confPath, CONF_FILE)
	conf.AddConfigPath(confPath)
	conf.SetConfigFile(confFile)

	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		if err := os.MkdirAll(confPath, 0755); err != nil {
			return fmt.Errorf("error creating confPath directory: %w", err)
		}
		file, err := os.Create(confFile)
		if err != nil {
			return fmt.Errorf("error creating config file: %w", err)
		}
		defer file.Close()
		_, err = file.WriteString("{}")
		if err != nil {
			return fmt.Errorf("error writing {} to file: %w", err)
		}
	}

	err = conf.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return fmt.Errorf("error reading config: %w", err)
		}
		err = conf.WriteConfigAs(confFile)
		if err != nil {
			return fmt.Errorf("error writing config file: %w", err)
		}
	}
	return nil
}

func setDefaults(c *viper.Viper) {
	c.SetDefault(BACKEND_URL_KEY, BACKEND_URL_DEFAULT)
	c.SetDefault(USERNAME_KEY, "")
	c.SetDefault(PASSWORD_KEY, "")
}

func GetConf() Config {
	return Config{
		BackendUrl: conf.GetString(BACKEND_URL_KEY),
		Username:   conf.GetString(USERNAME_KEY),
		Password:   conf.GetString(PASSWORD_KEY),
	}
}

func SetConf(c Config) error {
	conf.Set(BACKEND_URL_KEY, c.BackendUrl)
	conf.Set(USERNAME_KEY, c.Username)
	conf.Set(PASSWORD_KEY, c.Password)
	err := conf.WriteConfig()
	if err != nil {
		return fmt.Errorf("error updating conf: %w", err)
	}
	return nil
}

func GetBackendUrl() string {
	return conf.GetString(BACKEND_URL_KEY)
}
