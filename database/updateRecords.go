package database

import "database/sql"

func UpdateRecords(db *sql.DB, updateStatement string, id []uint8) (int64, error) {
	res, err := db.Exec(updateStatement, string(id))
	if err != nil {
		return 0, err
	}

	numUpdated, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return numUpdated, nil
}
