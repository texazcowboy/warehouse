package item

import (
	"database/sql"

	"github.com/texazcowboy/warehouse/internal/foundation/database"
)

func Create(item *Item, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		row := tx.QueryRow("INSERT INTO item(name) VALUES($1) RETURNING id", item.Name)
		if err := row.Scan(&item.ID); err != nil {
			return err
		}
		return nil
	})
}

func Get(id int64, db *sql.DB) (*Item, error) {
	var item Item
	err := database.ExecInTransaction(db, func(tx database.Transaction) error {
		row := tx.QueryRow("SELECT * FROM item WHERE id =$1", id)
		if err := row.Scan(&item.ID, &item.Name); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// pagination tbd.
func GetAll(db *sql.DB) ([]*Item, error) {
	var items []*Item
	err := database.ExecInTransaction(db, func(tx database.Transaction) error {
		rows, err := tx.Query("SELECT * from item") // nolint
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var i Item
			if err := rows.Scan(&i.ID, &i.Name); err != nil {
				return err
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

func Update(item *Item, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		res, err := tx.Exec("UPDATE item SET name=$1 WHERE id=$2", item.Name, item.ID)
		if err != nil {
			return err
		}
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if rowsAffected == 0 {
			return sql.ErrNoRows
		}
		return nil
	})
}

func Delete(id int64, db *sql.DB) error {
	return database.ExecInTransaction(db, func(tx database.Transaction) error {
		_, err := tx.Exec("DELETE FROM item WHERE id=$1", id)
		if err != nil {
			return err
		}
		return nil
	})
}
