package main

import (
	"database/sql"
	"fmt"
	"log"

	//Postgres driver
	_ "github.com/lib/pq"
)

func main() {
	// c := cron.New()
	// c.AddFunc("@every 1s", func() { fmt.Println("Every second") })
	// c.Start()
	// fmt.Println(c.Entries())
	//
	// time.Sleep(10 * time.Second)
	//
	// c.Stop()

	sqlStatement := `
    SELECT * FROM users
    WHERE new_user = TRUE;
    `
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", "34.71.218.171", 5432, "postgres",
		"***REMOVED***", "notify")
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Printf("%v", err)
	}

	rows, err := db.Query(sqlStatement)
	if err != nil {
		fmt.Println(err)
	}

	defer rows.Close()
	for rows.Next() {

		var id []uint8
		var firstName, lastName, email, username string
		var timeHr, timeMin float64
		var newUser bool

		err = rows.Scan(&id, &firstName, &lastName, &email, &username, &newUser, &timeHr, &timeMin)
		if err != nil {
			// handle this error
			fmt.Println(err)
		}
		fmt.Println("ID:", string(id))
		fmt.Println("Name:", firstName, lastName)
		fmt.Println("Email:", email)
		fmt.Println("Username:", username)
		fmt.Println("New user:", newUser)
		fmt.Printf("Time: %v:%v\n", timeHr, timeMin)

		updateStatement := `
    UPDATE users
    SET new_user = FALSE
    WHERE id = $1;
    `

		res, err := db.Exec(updateStatement, string(id))
		if err != nil {
			fmt.Println(err)
		}

		numUpdated, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Number of records updated:", numUpdated)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
	}

}
