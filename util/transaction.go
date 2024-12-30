package util

import "database/sql"

func InTx(db *sql.DB, callback func(*sql.Tx) error) func() error {
	return func() error {
		tx, err := db.Begin()
		if err != nil {
			return err
		}
		defer tx.Rollback()
		err = callback(tx)
		if err != nil {
			return err
		}
		err = tx.Commit()
		if err != nil {
			return err
		}
		return err
	}
}

func InTypedTx[K any, PK *K](db *sql.DB, callback func(*sql.Tx) (PK, error)) func() (PK, error) {
	return func() (PK, error) {
		tx, err := db.Begin()
		if err != nil {
			return nil, err
		}
		defer tx.Rollback()
		k, err := callback(tx)
		if err != nil {
			return nil, err
		}
		err = tx.Commit()
		if err != nil {
			return nil, err
		}
		return k, err
	}
}
