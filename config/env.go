package config

import (
	postgres "recro_demo/postgres"
)

// Env : Wrapper to hold all resources
type Env struct {
	DB      *postgres.DB
	Verbose string
	Config  *Config
}
