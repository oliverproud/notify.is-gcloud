package database

import "fmt"

//InsertUser inserts a new user into the database table 'users'
func InsertUser(firstName, lastName, email, username string) (string, error) {

	sqlStatement := `
  INSERT INTO users (first_name, last_name, email, username)
  VALUES ($1, $2, $3, $4)
  RETURNING id`

	var id []uint8
	err := db.QueryRow(sqlStatement, firstName, lastName, email, username).Scan(&id)
	if err != nil {
		return "", err
	}

	result := fmt.Sprintf("User inserted into DB")

	return result, nil
}