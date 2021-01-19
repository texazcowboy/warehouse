package user

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/texazcowboy/warehouse/internal/foundation/crypto"
	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

// Predefined queries for building statements.
const (
	createQuery        = "INSERT INTO usr(username, password) VALUES($1, $2) RETURNING id;"
	getByUsernameQuery = "SELECT * FROM usr WHERE username=$1;"
	deleteByIDQuery    = "DELETE FROM usr WHERE id=$1;"
)

type RepositoryInterface interface {
	Create(user *User) error
	GetByUsername(username string) (*User, error)
	DeleteByID(id int64) error
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{db}
}

type Repository struct {
	*sql.DB
}

func (r *Repository) Create(user *User) error {
	return database.ExecInTransaction(r.DB, func(tx database.Transaction) error {

		passwordHash, err := crypto.HashValue(user.Password)
		if err != nil {
			return errors.Wrap(err, "user: Create -> crypto.HashValue(***)")
		}

		stmt, err := tx.Prepare(createQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		row := stmt.QueryRow(user.Username, passwordHash)
		if err := row.Scan(&user.ID); err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	var user User
	err := database.ExecInTransaction(r.DB, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(getByUsernameQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		row := stmt.QueryRow(username)
		if err := row.Scan(&user.ID, &user.Username, &user.Password); err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) DeleteByID(id int64) error {
	return database.ExecInTransaction(r.DB, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(deleteByIDQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		_, err = stmt.Exec(id)
		if err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})
}
