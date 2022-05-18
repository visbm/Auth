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
	ErrUsernameExist  = errors.New("username already exist")
)

func (s *Store) NewDB(source string, logger *logging.Logger) error {
	db, err := sql.Open("sqlite3", "../../"+source)
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

	return nil

}
