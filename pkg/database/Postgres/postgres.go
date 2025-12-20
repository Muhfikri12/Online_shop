package postgres

import (
	"app/pkg/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Postgres struct {
	config config.DBConfig
}

func NewPostgres(config config.DBConfig) *Postgres {
	return &Postgres{
		config: config,
	}
}

func (p *Postgres) Connect() (*gorm.DB, error) {

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		p.config.DBHost, p.config.DBUser, p.config.DBPassword, p.config.DBName, p.config.DBPort,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	return db, nil
}
