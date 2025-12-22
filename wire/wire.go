package wire

import (
	"app/internal/route"
	"app/pkg/config"
	postgres "app/pkg/database/Postgres"
	"app/pkg/database/seeder"
	"log"

	"github.com/gin-gonic/gin"
)

func Initialize(cfg *config.Config) *gin.Engine {
	db, err := postgres.NewPostgres(cfg).Connect()
	if err != nil {
		panic(err)
	}

	router := route.Route(db, cfg)

	// seed data
	if cfg.DBConfig.Seed {
		log.Println("Seeding data...")
		seeder.SeedAll(db)
	}

	return router
}
