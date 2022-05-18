package app

import (
	"auth/internal/composite"
	"auth/internal/store/sqllite"
	"auth/pkg/logging"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	logger := logging.GetLogger()

	logger.Info("Initializing config...")
	config := Config{}
	config.NewConfig()

	sqllite := sqllite.Store{}

	logger.Info("Connecting to database via ", config.SQLLite3DataSourceName())
	err := sqllite.NewDB(config.SQLLite3DataSourceName(), &logger)
	if err != nil {
		logger.Fatal("Database err ", err)
	}
	defer sqllite.DB.Close()

	logger.Info("Initializing postgres composites...")
	var composite composite.Composites
	composite.NewComposites(sqllite)

	logger.Info("Initializing httprouter...")
	router := httprouter.New()
	ConfigureRouter(router, composite)

	logger.Info("Initializing server...")
	var server Server
	server.New(&config, router, &logger)
	if err := server.Start(); err != nil {
		logger.Fatal("Server falls")
	}
}
