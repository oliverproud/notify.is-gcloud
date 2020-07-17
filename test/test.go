package main

import (
	"database/sql"
	"fmt"
	"time"

	//Postgres driver
	_ "github.com/lib/pq"
)

type Args struct {
	t   time.Time
	lim int
}

func main() {

	sqlStatement := `SELECT id, first_name, email, username, timestamp FROM users WHERE EXTRACT(EPOCH FROM ($1 - timestamp)) > $2`

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", "34.71.218.171", 5432, "postgres",
		"***REMOVED***", "notify")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("Returning...")
		return
	}

	args := new(Args)
	args.t = time.Now()
	args.lim = 43200 // 12 hours in seconds

	rows, err := db.Query(sqlStatement, args.t, args.lim)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Returning...")
		return
	}

	defer rows.Close()
	for rows.Next() {

		var id []uint8
		var firstName, email, username string
		var timestamp time.Time

		err = rows.Scan(&id, &firstName, &email, &username,
			&timestamp)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}

		timeDiff := time.Since(timestamp)
		fmt.Printf("\nTime difference: %v\n", timeDiff)

		limit := time.Hour * 12

		if timeDiff > limit {
			fmt.Println("Time is greater than allowed")
			fmt.Println()
		} else {
			fmt.Println("Time OK")
			fmt.Println()
		}

		fmt.Println("ID:", string(id))
		fmt.Println("Name:", firstName)
		fmt.Println("Email:", email)
		fmt.Println("Username:", username)
		fmt.Printf("Timestamp: %v\n", timestamp)

		updateStatement := `
    UPDATE users
    SET timestamp = $1
    WHERE id = $2;
    `

		res, err := db.Exec(updateStatement, time.Now(), string(id))
		if err != nil {
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}

		numUpdated, err := res.RowsAffected()
		if err != nil {
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}
		fmt.Println("Number of records updated:", numUpdated)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		fmt.Println(err)
		fmt.Println("Returning...")
		return
	}
}
