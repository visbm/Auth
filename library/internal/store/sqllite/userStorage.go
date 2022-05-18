package sqllite

import (
	"auth/domain"
	"auth/internal/store"
	"auth/pkg/logging"
	"database/sql"
)

type userStorage struct {
	logger *logging.Logger
	db     *sql.DB
}

func NewUserStorage(db *sql.DB, logger *logging.Logger) store.UserStorage {
	return &userStorage{
		logger: logger,
		db:     db,
	}

}

func (us *userStorage) GetByUsername(username string) (*domain.Login, error) {
	var login domain.Login
	if err := us.db.QueryRow("SELECT * FROM users WHERE username = $1",
		username).Scan(
		&login.UUID,
		&login.Username,
		&login.Password,
	); err != nil {
		us.logger.Errorf("error occurred while selecting users from DB. err: %v", err)
		return nil, err
	}

	return &login, nil
}

func (us *userStorage) GetByUUID(UUID string) (*domain.Login, error) {
	var login domain.Login
	if err := us.db.QueryRow("SELECT * FROM users WHERE uuid = $1",
		UUID).Scan(
		&login.UUID,
		&login.Username,
		&login.Password,
	); err != nil {
		us.logger.Errorf("error occurred while selecting users from DB. err: %v", err)
		return nil, err
	}

	return &login, nil
}

func (us *userStorage) GetAll(limit, offset int) ([]*domain.Login, error) {
	rows, err := us.db.Query("SELECT * FROM users")
	if err != nil {
		us.logger.Errorf("error occurred while selecting all userss. err: %v", err)
		return nil, err
	}
	var logins []*domain.Login

	for rows.Next() {
		login := domain.Login{}
		err := rows.Scan(
			&login.UUID,
			&login.Username,
			&login.Password,
		)
		if err != nil {
			us.logger.Errorf("error occurred while selecting users. err: %v", err)
			continue
		}
		logins = append(logins, &login)
	}
	return logins, nil
}

func (us *userStorage) Create(login *domain.Login) (string, error) {

	var UUID string
	hashPass, err := domain.EncryptPassword(login.Password)
	if err != nil {
		us.logger.Errorf("error occurred while encrypting password . err: %v", err)
		return UUID, err
	}

	exist, err := us.CheckUsername(login.Username)
	if err != nil {
		us.logger.Errorf("error occurred while creating user. err: %v", err)
		return UUID, err
	}

	if *exist {
		us.logger.Errorf("username already exist")
		return UUID, ErrUsernameExist
	}

	if err = us.db.QueryRow(
		`INSERT INTO users (
                    username,
					password
		) VALUES ($1 , $2) RETURNING id`,
		login.Username,
		hashPass,
	).Scan(
		&UUID,
	); err != nil {
		us.logger.Errorf("error occurred while creating user. err: %v", err)
		return UUID, err
	}
	us.logger.Info("%v", UUID)

	return UUID, nil
}

func (us *userStorage) Delete(UUID string) error {
	result, err := us.db.Exec("DELETE FROM users WHERE uuid = $1", UUID)
	if err != nil {
		us.logger.Errorf("error occurred while deleting user. err: %v.", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		us.logger.Errorf("error occurred while deleting user (getting affected rows). err: %v", err)
		return err
	}

	if rowsAffected < 1 {
		us.logger.Errorf("error occurred while deleting user. err: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}
	us.logger.Infof("user with uuid %s was deleted.", UUID)
	return nil
}

func (us *userStorage) Update(login *domain.Login) error {
	result, err := us.db.Exec(
		`UPDATE users SET
	              full_name = COALESCE(NULLIF($1, ''), full_name)
		WHERE uuid = $2`,
		login.Username,
		login.UUID)

	if err != nil {
		us.logger.Errorf("error occurred while updating user. err: %v", err)
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		us.logger.Errorf("Error occurred while updating user. Err msg: %v.", err)
		return err
	}

	if rowsAffected < 1 {
		us.logger.Errorf("Error occurred while updating user. Err msg: %v.", ErrNoRowsAffected)
		return ErrNoRowsAffected
	}

	return nil
}

//Checks if the username exists

func (us *userStorage) CheckUsername(username string) (*bool, error) {
	var idIsExist bool
	err := us.db.QueryRow("SELECT EXISTS (SELECT username FROM users WHERE username = $1)", username).Scan(&idIsExist)
	if err != nil {
		us.logger.Errorf("Error occured while username checking. Err msg: %v", err)
		return &idIsExist, err
	}
	return &idIsExist, nil
}
