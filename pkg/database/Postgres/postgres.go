package postgres

import (
	"app/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Postgres struct {
	config *config.Config
}

func NewPostgres(config *config.Config) *Postgres {
	return &Postgres{
		config: config,
	}
}

func (p *Postgres) Connect() (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		p.config.DBConfig.DBHost, p.config.DBConfig.DBUser, p.config.DBConfig.DBPassword, p.config.DBConfig.DBName, p.config.DBConfig.DBPort,
	)

	logMode := logger.Silent
	if p.config.Environment == "local" || p.config.Environment == "development" {
		logMode = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logMode),
	})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
