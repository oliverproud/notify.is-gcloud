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
	"notify.is-go/sendgrid"
	"notify.is-go/timeDiff"
	//Postgres driver
	_ "github.com/lib/pq"
)

var id []uint8
var timestamp time.Time
var instagram, twitter bool
var firstName, email, username, updateStatement string

func runCheck() error {
	log.Println("Starting check...")

	selectStatement := `SELECT id, first_name, email, username, instagram, twitter, timestamp FROM users WHERE EXTRACT(EPOCH FROM ((now() at time zone 'utc') - timestamp)) > 43200.0`

	rows, err := database.SelectRecords(db, selectStatement)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {

		// Default update statement
		updateStatement = `
				UPDATE users
				SET timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`

		if err = rows.Scan(&id, &firstName, &email, &username, &instagram, &twitter, &timestamp); err != nil {
			return err
		}

		timeDiff.CalculateDiff(timestamp)

		fmt.Println("ID:", string(id))
		fmt.Println("Name:", firstName)
		fmt.Println("Email:", email)
		fmt.Println("Username:", username)
		fmt.Println("Instagram:", instagram)
		fmt.Println("Twitter:", twitter)

		if instagram && twitter {
			fmt.Println("Run both Instagram and Twiter")
			if err := check.RunHeadless(email, firstName, username); err != nil {
				return err
			}
			if check.Available {
				resp, err := sendgrid.SuccessEmailInstagram(email, firstName, username)
				if err != nil {
					return err
				}
				fmt.Println("Sendgrid Response:", resp.StatusCode)
			}

			available, err := check.TwitterAPI(username)
			if err != nil {
				return err
			}

			// Update statement
			// Send success email
			if available {
				resp, err := sendgrid.SuccessEmailTwitter(email, firstName, username)
				if err != nil {
					return err
				}
				fmt.Println("Sendgrid Response:", resp.StatusCode)
			}
			if check.Available && available {
				updateStatement = `
				UPDATE users
				SET instagram = false, twitter = false, timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`
			} else if check.Available {
				updateStatement = `
				UPDATE users
				SET instagram = false, timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`
			} else if available {
				updateStatement = `
				UPDATE users
				SET twitter = false, timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`
			}
		} else if instagram {
			fmt.Println("Only run Instagram")
			if err := check.RunHeadless(email, firstName, username); err != nil {
				return err
			}
			if check.Available {
				resp, err := sendgrid.SuccessEmailInstagram(email, firstName, username)
				if err != nil {
					return err
				}
				fmt.Println("Sendgrid Response:", resp.StatusCode)

				updateStatement = `
				UPDATE users
				SET instagram = false, timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`
			}
		} else if twitter {
			fmt.Println("Only run Twitter")
			available, err := check.TwitterAPI(username)
			if err != nil {
				return err
			}

			// Update statement
			// Send success email
			if available {
				resp, err := sendgrid.SuccessEmailTwitter(email, firstName, username)
				if err != nil {
					return err
				}
				fmt.Println("Sendgrid Response:", resp.StatusCode)
				updateStatement = `
				UPDATE users
				SET twitter = false, timestamp = (now() at time zone 'utc')
				WHERE id = $1;
				`
			}
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

	// Setenv here
	
	const (
		port   = 5432
		user   = "postgres"
		dbName = "notify"
	)

	var host = os.Getenv("DB_HOST")
	var password = os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbName)

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
