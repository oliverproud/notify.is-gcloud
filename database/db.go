package database

import (
	"database/sql"
	"log"

	//Postgres driver
	_ "github.com/lib/pq"
)

var db *sql.DB

//InitDB initialises a database instance
func InitDB(dataSourceName string) {
	var err error
	db, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Printf("%v", err)
	}

	if err = db.Ping(); err != nil {
		log.Printf("%v", err)
	}
}
