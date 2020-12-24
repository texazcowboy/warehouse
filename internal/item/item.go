package item

import "database/sql"

func create(name string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec("INSERT INTO item(name) VALUES($1)", name)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func get(id int64, db *sql.DB) (*Item, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	var item Item
	row := db.QueryRow("SELECT * FROM item WHERE id =$1", id)
	if err = row.Scan(&item.Id, &item.Name); err != nil {
		tx.Rollback()
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return &item, err
}

func getAll(db *sql.DB) ([]*Item, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	var items []*Item

	rows, err := db.Query("SELECT * from item")
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var i Item
		if err := rows.Scan(&i.Id, &i.Name); err != nil {
			tx.Rollback()
			return nil, err
		}
		items = append(items, &i)
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return items, nil
}

func update(id int64, name string, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec("UPDATE item SET name=$1 WHERE id=$2", name, id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func delete(id int64, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	_, err = db.Exec("DELETE FROM item WHERE id=$1", id)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}
