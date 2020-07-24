package database

import "database/sql"

func SelectRecords(db *sql.DB, selectStatement string) (*sql.Rows, error) {

	rows, err := db.Query(selectStatement)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
