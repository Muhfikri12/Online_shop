package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Debug       bool
	Timezone    string
	Port        string
	Location    *time.Location `anonymous:"true"`
	Name        string
	Environment string
	PrivateKey  string
	PublicKey   string
	DBConfig    DBConfig
}

type DBConfig struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	Seed       bool
	Migrate    bool
}

func NewConfig() *Config {
	viper.SetDefault("TIMEZONE", "Asia/Jakarta")
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("APP_NAME", "Online Shop")
	viper.SetDefault("APP_ENV", "local")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "postgres")

	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	_ = viper.ReadInConfig()

	return &Config{
		Debug:       viper.GetBool("DEBUG"),
		Timezone:    viper.GetString("TIMEZONE"),
		Port:        viper.GetString("PORT"),
		Location:    time.Local,
		Name:        viper.GetString("APP_NAME"),
		Environment: viper.GetString("APP_ENV"),
		PrivateKey:  viper.GetString("PRIVATE_KEY"),
		PublicKey:   viper.GetString("PUBLIC_KEY"),
		DBConfig: DBConfig{
			DBHost:     viper.GetString("DB_HOST"),
			DBPort:     viper.GetString("DB_PORT"),
			DBUser:     viper.GetString("DB_USER"),
			DBPassword: viper.GetString("DB_PASSWORD"),
			DBName:     viper.GetString("DB_NAME"),
			Seed:       viper.GetBool("DB_SEED"),
			Migrate:    viper.GetBool("DB_MIGRATE"),
		},
	}
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	Prefix   string
}

func (c *Config) RedisConfig() *RedisConfig {
	return &RedisConfig{
		Host:     viper.GetString("REDIS_HOST"),
		Port:     viper.GetString("REDIS_PORT"),
		Password: viper.GetString("REDIS_PASSWORD"),
		Prefix:   viper.GetString("REDIS_PREFIX"),
	}
}
