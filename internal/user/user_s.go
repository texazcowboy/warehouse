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
func Delete(id int64, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		_, err := tx.Exec("DELETE FROM usr WHERE id=$1", id)
		if err != nil {
			return errors.Wrapf(err, "user id: %v", id)
		}
		return nil
	})
}
