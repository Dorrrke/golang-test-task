package config

import (
	"flag"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env          string        `yaml:"env" env-default:"local"`
	Storage      StorageConfig `yaml:"storage"`
	ServerConfig ServerConfig  `yaml:"server"`
}

type StorageConfig struct {
	User   string `yaml:"user"`
	Pass   string `yaml:"password"`
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DbName string `yaml:"dbname"`
}

type ServerConfig struct {
	Host    string        `yaml:"host"`
	Port    string        `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file is not exist")
	}

	var config Config
	if err := cleanenv.ReadConfig(path, &config); err != nil {
		panic("failed to read config: " + err.Error())
	}

	return &config
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()
	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func (serverCfg *ServerConfig) GetServerAddr() string {
	if serverCfg.Host == "localhost" {
		return ":" + serverCfg.Port
	}
	return serverCfg.Host + ":" + serverCfg.Port
}
