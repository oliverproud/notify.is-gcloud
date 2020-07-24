package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"notify.is-go/check"
	"notify.is-go/database"
	"notify.is-go/timeDiff"
	//Postgres driver
	_ "github.com/lib/pq"
)

var id []uint8
var firstName, email, username, updateStatement string
var timestamp time.Time

func runCheck() error {
	log.Println("Starting check...")

	selectStatement := `SELECT id, first_name, email, username, timestamp FROM users WHERE EXTRACT(EPOCH FROM ((now() at time zone 'utc') - timestamp)) > 43200.0`

	rows, err := database.SelectRecords(db, selectStatement)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {

		if err = rows.Scan(&id, &firstName, &email, &username, &timestamp); err != nil {
			return err
		}

		timeDiff.CalculateDiff(timestamp)

		fmt.Println("ID:", string(id))
		fmt.Println("Name:", firstName)
		fmt.Println("Email:", email)
		fmt.Println("Username:", username)

		if err = check.RunHeadless(email, firstName, username); err != nil {
			return err
		}

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

		numUpdated, err := database.UpdateRecords(db, updateStatement, id)
		if err != nil {
			return err
		}
		fmt.Println("Number of records updated:", numUpdated)
	}

	// get any error encountered during iteration
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["auth"]
	if !ok || len(keys[0]) < 1 || keys[0] != os.Getenv("SERVER_PASSWORD") {
		fmt.Fprintf(w, "You are not authorised to access this page")
	} else {
		log.Print("Notify.is: received a request")

		if err := runCheck(); err != nil {
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}
		fmt.Fprintf(w, "Ready to process requests.\n")
	}
}

var db *sql.DB

func init() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", os.Getenv("DB_HOST"), 5432, "postgres", os.Getenv("DB_PASSWORD"), "notify")

	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		fmt.Printf("%v", err)
		fmt.Println("Returning...")
		return
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	log.Print("Notify.is: starting server...")

	http.HandleFunc("/", handler)

	log.Printf("Notify.is: listening on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
