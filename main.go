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
