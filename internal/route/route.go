package route

import (
	_ "app/docs" // swagger docs (generate with: swag init -g main.go -o docs)
	"app/internal/app/handler"
	"app/internal/app/repository"
	"app/internal/app/service"
	"app/pkg/config"
	rds "app/pkg/database/redis"
	"app/pkg/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/gorm"
)

func Route(db *gorm.DB, cfg *config.Config, rds rds.Redis) *gin.Engine {

	handler, router := InitRoute(db, cfg, rds)

	api := router.Group("/api")
	api.Use(middleware.Meta())

	// public auth routes
	api.POST("/login", handler.User.Login)
	api.POST("/refresh-token", handler.User.RefreshToken)

	// protected routes
	api.Use(middleware.Auth(rds, cfg))

	api.POST("/logout", handler.User.Logout)

	// Product
	api.GET("/product/:uuid", handler.Product.FindByUUID)
	api.GET("/products", handler.Product.FindAll)

	return router

}

func InitRoute(db *gorm.DB, cfg *config.Config, rds rds.Redis) (*handler.Handler, *gin.Engine) {
	// Repository
	repo := repository.NewRepository(db)

	// Service
	service := service.NewService(repo, cfg, rds)

	// Handler
	handler := handler.NewHandler(service, cfg)

	// Route
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Swagger UI (no auth)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	ginMode := gin.ReleaseMode
	if ginMode == gin.DebugMode {
		gin.SetMode(gin.DebugMode)
	}

	return handler, router

}
