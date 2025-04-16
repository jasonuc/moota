package store

import "database/sql"

type Transaction struct {
	tx *sql.Tx
}

func (t *Transaction) Rollback() error {
	return t.tx.Rollback()
}

func (t *Transaction) Commit() error {
	return t.tx.Commit()
}

type dbOrTx interface {
	Query(string, ...any) (*sql.Rows, error)
	Exec(string, ...any) (sql.Result, error)
	QueryRow(string, ...any) *sql.Row
}
