package config

import (
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

type ServerSettings struct {
	ListenAddress string `yaml:"listen_address"`
}

type AuthSettings struct {
	ApiToken string `yaml:"api_token"`
	ApiLogin string `yaml:"api_login"`
}

type DatabaseSettings struct {
	MgoDsn         string        `yaml:"mgo_dsn"`
	Retries        int           `yaml:"retries"`
	ConnectTimeout time.Duration `yaml:"connect_timeout"`
}

type Config struct {
	Server   ServerSettings   `yaml:"server_settings"`
	Database DatabaseSettings `yaml:"database_settings"`
	Auth     AuthSettings     `yaml:"auth_settings"`
}

var (
	defaultConfig *Config
	BuildCommit   string
	BuildDate     string
)

func Load(yamlPath string) {
	configFile, err := os.Open(yamlPath)
	if err != nil {
		panic(err)
	}

	data, err := ioutil.ReadAll(configFile)
	if err != nil {
		panic(err)
	}

	defaultConfig = &Config{}
	if err := yaml.Unmarshal(data, defaultConfig); err != nil {
		panic(err)
	}
}

func SetConfig(config *Config) {
	defaultConfig = config
}

func GetConfig() Config {
	return *defaultConfig
}
