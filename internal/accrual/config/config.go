package config

import "time"

type HTTPConfig struct {
	RunAddress string
}

type PostgresConfig struct {
	DatabaseURI    string
	ConnectTimeout time.Duration
}
