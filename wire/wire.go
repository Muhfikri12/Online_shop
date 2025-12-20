package wire

import (
	"app/internal/route"
	"app/pkg/config"
	postgres "app/pkg/database/Postgres"

	"github.com/gin-gonic/gin"
)

func Initialize() *gin.Engine {
	db, err := postgres.NewPostgres(config.DBConfig{}).Connect()
	if err != nil {
		panic(err)
	}

	router := route.Route(db)

	return router
}
