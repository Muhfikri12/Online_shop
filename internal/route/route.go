package route

import (
	"app/internal/app/handler"
	"app/internal/app/repository"
	"app/internal/app/service"
	"app/pkg/config"
	"app/pkg/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func Route(db *gorm.DB, cfg *config.Config) *gin.Engine {

	handler, router := InitRoute(db, cfg)

	api := router.Group("/api")
	api.Use(middleware.Meta())

	api.POST("/login", handler.User.Login)

	// Product
	api.GET("/product/:uuid", handler.Product.FindByUUID)

	return router

}

func InitRoute(db *gorm.DB, cfg *config.Config) (*handler.Handler, *gin.Engine) {
	// Repository
	repo := repository.NewRepository(db)

	// Service
	service := service.NewService(repo, cfg)

	// Handler
	handler := handler.NewHandler(service)

	// Route
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	ginMode := gin.ReleaseMode
	if ginMode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	return handler, router

}
