package main

import (
	"fmt"
	"graves/internal/repository"
	routers "graves/internal/router"
	"graves/pkg/config"
	"log"
)

func main() {
	// This is the main entry point for the application.
	// The actual implementation will be in the cmd package.

	cfg, err := config.GetInstance()
	if err != nil || cfg.Server == nil {
		log.Fatal(fmt.Sprintf("Error getting config instance: %v", err))
	}

	if _, err := repository.GetInstance(); err != nil {
		log.Fatal(fmt.Sprintf("Error getting repository instance: %v", err))
	}

	r := routers.SetupRouter()
	r.Run(fmt.Sprintf("%v:%v", cfg.Server.Host, cfg.Server.Port))
}
