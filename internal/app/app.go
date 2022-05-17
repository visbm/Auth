package app

import (
	"auth/internal/composite"
	"auth/internal/store/sqllite"
	"auth/pkg/logging"
	"database/sql"
	"fmt"

	"github.com/julienschmidt/httprouter"
	_ "github.com/mattn/go-sqlite3"
)

func Run() {
	logger := logging.GetLogger()

	logger.Info("Initializing config...")
	config := Config{}
	config.NewConfig()

	sqllite := sqllite.Store{}

	logger.Info("Connecting to database via %s", config.SQLLite3DataSourceName())
	db, err := sql.Open("sqlite3", "file:test.db?cache=shared&mode=memory")
	if err != nil {
		logger.Fatal("Database falls", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		logger.Fatal("Database falls", err)
	}
	sqllite.NewDB(db, &logger)

	logger.Info("Initializing postgres composites...")
	var composite composite.Composites
	composite.NewComposites(sqllite)

	s, err := db.Prepare(`CREATE TABLE products(
		id INTEGER PRIMARY KEY AUTOINCREMENT, 
		model TEXT,
		company TEXT,
		price INTEGER
	  );`)
	if err != nil {
		logger.Fatal("Database err", err)
	}
	s.Exec()

	result, err := db.Exec("insert into products (model, company, price) values ('iPhone X', $1, $2)",
		"Apple", 72000)
	if err != nil {
		logger.Error("Database falls", err)
	}

	fmt.Println(result.LastInsertId()) // id последнего добавленного объекта
	fmt.Println(result.RowsAffected()) // количество добавленных строк

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
