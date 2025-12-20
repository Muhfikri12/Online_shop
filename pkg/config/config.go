package config

import "time"

type Config struct {
	Debug       bool
	Timezone    string
	Port        string
	Location    *time.Location `anonymous:"true"`
	Name        string
	Environment string
	DBConfig    DBConfig
}

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}
