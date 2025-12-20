package route

import (
	"app/internal/app/handler"
	"app/internal/app/repository"
	"app/internal/app/service"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(db *gorm.DB) *gin.Engine {

	handler := InitRoute(db)

	// Route
	router := gin.New()
	ginMode := gin.ReleaseMode
	if ginMode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	api := router.Group("/api")
	api.GET("/product/:uuid", handler.Product.FindByUUID)

	return router

}

func InitRoute(db *gorm.DB) *handler.Handler {
	// Repository
	repo := repository.NewRepository(db)

	// Service
	service := service.NewService(repo)

	// Handler
	return handler.NewHandler(service)
}
