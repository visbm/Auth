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
)

func (s *Store) NewDB(db *sql.DB, logger *logging.Logger) {
	s.Logger = logger
	s.DB = db
	s.UserStorage = NewUserStorage(s.DB, s.Logger)

}
