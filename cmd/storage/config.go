package main

import (
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Debug         bool   `envconfig:"DEBUG"`
	Port          string `envconfig:"PORT" required:"true"`
	DBConnString  string `envconfig:"POSTGRES_CONNSTR" required:"true"`
	MinioUser     string `envconfig:"MINIO_USER" required:"true"`
	MinioPassword string `envconfig:"MINIO_PASSWORD" required:"true"`
	MinioEndpoint string `envconfig:"MINIO_ENDPOINT" required:"true"`
	FOBucket      string `envconfig:"FO_BUCKET" required:"true"`
}

func ReadConfig() (Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}
