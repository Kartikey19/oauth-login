package config

import (
	postgres "recro_demo/postgres"
)

// Env : Wrapper to hold all resources
type Env struct {
	Db      *postgres.DB
	Verbose string
	Config  *Config
}
