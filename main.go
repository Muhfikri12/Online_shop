// @title           Online Shop API
// @version         1.0
// @description     REST API for Online Shop with JWT authentication, refresh tokens, and remember-me.
// @termsOfService  http://swagger.io/terms/
// @contact.name    API Support
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
// @Host            localhost:8080
// @BasePath        /api
// @schemes         http
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"app/cmd"
	"app/pkg/config"
	"app/wire"
)

func main() {
	// config
	cfg := config.NewConfig()

	// wire
	engine := wire.Wire(cfg)

	// run
	cmd.ApiServer(cfg.Port, cfg.Name, engine)
}
