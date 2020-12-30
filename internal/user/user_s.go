package user

import (
	"database/sql"

	"github.com/pkg/errors"

	"github.com/texazcowboy/warehouse/internal/foundation/crypto"

	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

func Create(user *User, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		passwordHash, err := crypto.HashValue(user.Password)
		if err != nil {
			return errors.Wrap(err, "Create -> crypto.HashValue(***)")
		}
		row := tx.QueryRow("INSERT INTO usr(username, password) VALUES($1, $2) RETURNING id", user.Username, passwordHash)
		if err := row.Scan(&user.ID); err != nil {
			return err
		}
		return nil
	})
}

func GetByUsername(username string, db *sql.DB) (*User, error) {
	var user User
	err := database.ExecInTransaction(db, func(tx database.Transaction) error {
		row := tx.QueryRow("SELECT * FROM usr WHERE username=$1", username)
		if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrap(err, "GetByUsername -> database.ExecInTransaction(...)")
	}
	return &user, nil
}

func Delete(id int64, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		_, err := tx.Exec("DELETE FROM usr WHERE id=$1", id)
		if err != nil {
			return errors.Wrapf(err, "user id: %v", id)
		}
		return nil
	})
}
