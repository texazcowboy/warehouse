package item

import (
	"database/sql"

	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

// Predefined queries for building statements.
const (
	createQuery     = "INSERT INTO item(name) VALUES($1) RETURNING id;"
	getByIDQuery    = "SELECT * FROM item WHERE id =$1;"
	getAllQuery     = "SELECT * from item;"
	updateQuery     = "UPDATE item SET name=$1 WHERE id=$2;"
	deleteByIDQuery = "DELETE FROM item WHERE id=$1;"
)

type RepositoryInterface interface {
	Create(item *Item) error
	GetByID(id int64) (*Item, error)
	GetAll() ([]*Item, error)
	Update(item *Item) (int64, error)
	DeleteByID(id int64) (int64, error)
}

func NewRepository(db *sql.DB) RepositoryInterface {
	return &Repository{db: db}
}

type Repository struct {
	db *sql.DB
}

func (r *Repository) Create(item *Item) error {
	return database.ExecInTransaction(r.db, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(createQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		row := stmt.QueryRow(item.Name)
		if err := row.Scan(&item.ID); err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})
}

func (r *Repository) GetByID(id int64) (*Item, error) {

	var item Item

	err := database.ExecInTransaction(r.db, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(getByIDQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		row := stmt.QueryRow(id)
		if err := row.Scan(&item.ID, &item.Name); err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &item, nil
}

func (r *Repository) GetAll() ([]*Item, error) {

	var items []*Item

	err := database.ExecInTransaction(r.db, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(getAllQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		rows, err := stmt.Query() // nolint
		if err != nil {
			return database.NewSQLError(err)
		}
		defer rows.Close()

		for rows.Next() {
			var i Item
			if err := rows.Scan(&i.ID, &i.Name); err != nil {
				return database.NewSQLError(err)
			}
			items = append(items, &i)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *Repository) Update(item *Item) (int64, error) {

	var rowsUpdated int64

	err := database.ExecInTransaction(r.db, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(updateQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(item.Name, item.ID)
		if err != nil {
			return database.NewSQLError(err)
		}

		rowsUpdated, err = res.RowsAffected()
		if err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})

	return rowsUpdated, err
}

func (r *Repository) DeleteByID(id int64) (int64, error) {

	var rowsDeleted int64

	err := database.ExecInTransaction(r.db, func(tx database.Transaction) error {

		stmt, err := tx.Prepare(deleteByIDQuery)
		if err != nil {
			return database.NewSQLError(err)
		}
		defer stmt.Close()

		res, err := stmt.Exec(id)
		if err != nil {
			return database.NewSQLError(err)
		}

		rowsDeleted, err = res.RowsAffected()
		if err != nil {
			return database.NewSQLError(err)
		}

		return nil
	})

	return rowsDeleted, err
}
