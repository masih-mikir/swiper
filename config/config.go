package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ruizu/gcfg"
)

type Config struct {
	Server   ServerConfig
	Account  AccountConfig
	Redis    RedisConfig
	InMemory InMemoryConfig
}

type ServerConfig struct {
	Enviroment string
	DBTimeout  time.Duration
}

type AccountConfig struct {
	Port     string
	MasterDB string
	SlaveDB  string
}

type RedisConfig struct {
	Host        string
	PoolSize    int
	DialTimeout time.Duration
	IdleTimeout time.Duration
}

type InMemoryConfig struct {
	DefaultExpiration time.Duration
	IntervalPurges    time.Duration
}

func InitConfig(configPaths ...string) (*Config, bool) {
	var cfg Config
	var ok bool

	env := os.Getenv("SYSENV")
	if env == "" {
		env = "development"
	}

	for _, configPath := range configPaths {
		configPath = fmt.Sprintf("%s/config.%s.ini", configPath, env)
		fmt.Println(configPath)
		if err := gcfg.ReadFileInto(&cfg, configPath); err != nil {
			log.Println(err)
		} else {
			ok = true
			log.Printf("open %s succcessful", configPath)
			break
		}
	}

	return &cfg, ok
}
