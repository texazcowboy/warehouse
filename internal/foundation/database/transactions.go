package database

import (
	"database/sql"
)

type Transaction interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
	Prepare(query string) (*sql.Stmt, error)
	Query(query string, args ...interface{}) (*sql.Rows, error)
	QueryRow(query string, args ...interface{}) *sql.Row
}

type TxFn func(Transaction) error

func ExecInTransaction(db *sql.DB, fn TxFn) (err error) {
	tx, err := db.Begin()
	if err != nil {
		return
	}
	defer func() {
		switch p := recover(); {
		case p != nil:
			tx.Rollback() // nolint
			panic(p)
		case err != nil:
			tx.Rollback() // nolint
		default:
			err = tx.Commit()
		}
	}()

	err = fn(tx)
	return
}
