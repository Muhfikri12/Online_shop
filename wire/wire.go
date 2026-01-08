package wire

import (
	"app/internal/route"
	"app/pkg/config"
	postgres "app/pkg/database/Postgres"
	"app/pkg/database/redis"
	"app/pkg/database/seeder"
	"log"

	"github.com/gin-gonic/gin"
)

func Wire(cfg *config.Config) *gin.Engine {
	db, err := postgres.NewPostgres(cfg).Connect()
	if err != nil {
		panic(err)
	}

	// redis config
	redisConfig := cfg.RedisConfig()

	// set Redis
	rds := redis.NewRedis(redisConfig, 0)

	router := route.Route(db, cfg, rds)

	// seed data
	if cfg.DBConfig.Seed {
		log.Println("Seeding data...")
		seeder.SeedAll(db)
	}

	return router
}
