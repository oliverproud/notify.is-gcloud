package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"notify.is-go/check"
	//Postgres driver
	_ "github.com/lib/pq"
)

var id []uint8
var firstName, email, username string
var timestamp time.Time

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

func selectUsers(db *sql.DB, selectStatement string) (*sql.Rows, error) {

	rows, err := db.Query(selectStatement)
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

func runThis() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", os.Getenv("DB_HOST"), 5432, "postgres",
		os.Getenv("DB_PASSWORD"), "notify")

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("Returning...")
		return
	}

	log.Println("Starting check...")

	selectStatement := `SELECT id, first_name, email, username, timestamp FROM users WHERE EXTRACT(EPOCH FROM ((now() at time zone 'utc') - timestamp)) > 43200.0`

	rows, err := selectUsers(db, selectStatement)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Returning...")
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

		err = check.RunCheck(email, firstName, username)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}

		var updateStatement string

		// For testing
		// updateStatement = `
		// 	UPDATE users
		// 	SET timestamp = (now() at time zone 'utc')
		// 	WHERE id = $1;
		// 	`

		if check.Available {
			updateStatement = `
		  DELETE FROM users
		  WHERE id = $1;
			`
		} else {
			updateStatement = `
		  UPDATE users
		  SET timestamp = (now() at time zone 'utc')
		  WHERE id = $1;
			`
		}

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
}

func handler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["auth"]
	if !ok || len(keys[0]) < 1 || keys[0] != os.Getenv("SERVER_PASSWORD") {
		fmt.Fprintf(w, "You are not authorised to access this page")
	} else {
		log.Print("Notify.is: received a request")

		runThis()

		fmt.Fprintf(w, "Ready to process requests.\n")
	}

}

func main() {
	log.Print("Notify.is: starting server...")

	http.HandleFunc("/", handler)

	log.Printf("Notify.is: listening on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
