package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug        bool   `envconfig:"DEBUG"`
	Port         string `envconfig:"PORT" required:"true"`
	DBConnString string `envconfig:"POSTGRES_CONNSTR" required:"true"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
