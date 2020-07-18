package main

import (
	"database/sql"
	"fmt"
	"os"
	"os/signal"
	"time"

	//Postgres driver
	_ "github.com/lib/pq"
	"github.com/robfig/cron"
)

var id []uint8
var firstName, email, username string
var timestamp time.Time

type Args struct {
	t   time.Time
	lim int
}

func timeDiff(timestamp time.Time) {
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
}

func selectUsers(db *sql.DB, args *Args, selectStatement string) (*sql.Rows, error) {

	rows, err := db.Query(selectStatement, args.t, args.lim)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func updateTimestamp(db *sql.DB, updateStatement string, id []uint8) (int64, error) {
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

func main() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", "34.71.218.171", 5432, "postgres",
		"***REMOVED***", "notify")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("Returning...")
		return
	}

	args := new(Args)
	args.t = time.Now().UTC()
	args.lim = 43200 // 12 hours in seconds

	c := cron.New()
	c.AddFunc("@every 1m", func() {

		fmt.Println("Starting check...")

		selectStatement := `SELECT id, first_name, email, username, timestamp FROM users WHERE EXTRACT(EPOCH FROM ($1 - timestamp)) > $2`

		rows, err := selectUsers(db, args, selectStatement)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Returning")
			return
		}

		defer rows.Close()
		for rows.Next() {

			err = rows.Scan(&id, &firstName, &email, &username, &timestamp)
			if err != nil {
				fmt.Println(err)
				fmt.Println("Returning...")
				return
			}

			timeDiff(timestamp)

			fmt.Println("ID:", string(id))
			fmt.Println("Name:", firstName)
			fmt.Println("Email:", email)
			fmt.Println("Username:", username)
			fmt.Printf("Timestamp: %v\n", timestamp)

			updateStatement := `
	    UPDATE users
	    SET timestamp = (now() at time zone 'utc')
	    WHERE id = $1;
    	`
			numUpdated, err := updateTimestamp(db, updateStatement, id)
			if err != nil {
				fmt.Println(err)
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
	})

	c.Start()

	// Wait for a signal to quit:
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan

}
