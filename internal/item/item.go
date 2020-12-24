package item

import "database/sql"

func Create(item *Item, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	row := db.QueryRow("INSERT INTO item(name) VALUES($1) RETURNING id", item.Name)
	if err = row.Scan(&item.Id); err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func Get(id int64, db *sql.DB) (*Item, error) {
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
	return &item, nil
}

// pagination tbd
func GetAll(db *sql.DB) ([]*Item, error) {
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

func Update(item *Item, db *sql.DB) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	res, err := db.Exec("UPDATE item SET name=$1 WHERE id=$2", item.Name, item.Id)
	if err != nil {
		tx.Rollback()
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		tx.Rollback()
		return err
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	return nil
}

func Delete(id int64, db *sql.DB) error {
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
