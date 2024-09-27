package user

import (
	"database/sql"
	"errors"
)

type UserStorage struct {
	db *sql.DB
}

func NewUserStorage(db *sql.DB) *UserStorage {
	return &UserStorage{
		db: db,
	}
}

func (us *UserStorage) RegisterNewUserStorage(user *UserPayload) error {
	query := "INSERT INTO users(username, email, password_hash) VALUES($1, $2, $3)"
	_, err := us.db.Exec(query, user.Name, user.Email, user.Password)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserStorage) VerifiyUserStorage(userpayload *UserPayload) (User, error) {
	var user User
	query := "SELECT * FROM users WHERE email=$1"
	result, err := us.db.Query(query, userpayload.Email)
	if err != nil {
		return User{}, err
	}

	if result.Next() {
		err = result.Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return User{}, err
		}
	} else {
		return User{}, errors.New("not Found")
	}

	return user, nil
}
