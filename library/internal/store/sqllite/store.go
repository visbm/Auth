package sqllite

import (
	"auth/internal/store"
	"auth/pkg/logging"
	"database/sql"
	"errors"
)

type Store struct {
	Logger      *logging.Logger
	DB          *sql.DB
	UserStorage store.UserStorage
}

var (
	// ErrNoRowsAffected ...
	ErrNoRowsAffected = errors.New("no rows affected")
	ErrUsernameExist = errors.New("username already exist")
)

func (s *Store) NewDB(source string, logger *logging.Logger) error {
	db, err := sql.Open("sqlite3", "../../database/sqlite3/data/db.sqlite")
	if err != nil {
		logger.Fatal("Database falls", err)
		return err
	}

	err = db.Ping()
	if err != nil {
		logger.Fatal("Database falls", err)
		return err

	}

	s.Logger = logger
	s.DB = db
	s.UserStorage = NewUserStorage(db, s.Logger)

	/*k, err := db.Prepare(`CREATE TABLE IF NOT EXISTS users(
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT,
		password TEXT
	  );`)
	if err != nil {
		logger.Error("Create table err ", err)
	}
	k.Exec()
	defer k.Close()

	_, err = k.Exec("INSERT INTO users (username, password) values ($1, $2)",
		"admin", "admin")
	if err != nil {
		logger.Error("insert table err ", err)
	}
	k.Close()*/

	return nil

}
