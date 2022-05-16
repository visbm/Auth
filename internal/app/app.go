package app

import (
	"auth/internal/composite"
	"auth/pkg/logging"

	"github.com/julienschmidt/httprouter"
)

func Run() {
	logger := logging.GetLogger()

	logger.Info("Initializing config...")
	config := Config{}


	config.NewConfig()


	logger.Info("Initializing postgres composites...")
	var authcomposite composite.AuthComposite
	authcomposite.New(&logger)

	logger.Info("Initializing httprouter...")
	router := httprouter.New()
	ConfigureRouter(router, authcomposite)

	logger.Info("Initializing server...")
	var server Server
	server.New(&config, router, &logger)
	if err := server.Start(); err != nil {
		logger.Fatal("Server falls")
	}
}
