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
	"notify.is-go/statements"
	"notify.is-go/timeDiff"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	//Postgres driver
	_ "github.com/lib/pq"
)

var id []uint8
var timestamp time.Time
var instagram, twitter, github bool
var firstName, email, username, updateStatement string

// Checks Instagram, sends email, updates database
func runInstagramCheck(email, firstName, username string) error {
	instagramAvailable, err := check.Instagram(username)
	if err != nil {
		return err
	}

	if instagramAvailable {
		resp, err := sendgrid.SuccessEmailInstagram(email, firstName, username)
		if err != nil {
			return err
		}
		numUpdated, err := database.UpdateRecords(db, statements.InstagramUpdateStatement, id)
		if err != nil {
			return err
		}
		fmt.Println("Sendgrid Response:", resp.StatusCode)
		fmt.Println("Number of records updated:", numUpdated)
	}
	return nil
}

// Checks Twitter, sends email, updates database
func runTwitterCheck(email, firstName, username string) error {
	twitterAvailable, err := check.Twitter(username)
	if err != nil {
		return err
	}

	if twitterAvailable {
		resp, err := sendgrid.SuccessEmailTwitter(email, firstName, username)
		if err != nil {
			return err
		}
		numUpdated, err := database.UpdateRecords(db, statements.TwitterUpdateStatement, id)
		if err != nil {
			return err
		}
		fmt.Println("Sendgrid Response:", resp.StatusCode)
		fmt.Println("Number of records updated:", numUpdated)
	}
	return nil
}

// Checks GitHub, sends email, updates database
func runGithubCheck(email, firstName, username string) error {

	githubAvailable, err := check.Github(username)
	if err != nil {
		return err
	}

	if githubAvailable {
		resp, err := sendgrid.SuccessEmailGithub(email, firstName, username)
		if err != nil {
			return err
		}
		numUpdated, err := database.UpdateRecords(db, statements.GithubUpdateStatement, id)
		if err != nil {
			return err
		}
		fmt.Println("Sendgrid Response:", resp.StatusCode)
		fmt.Println("Number of records updated:", numUpdated)
	}
	return nil
}

// Selects records from database runs checks on them
func runCheck() error {
	log.Println("Starting check...")

	rows, err := database.SelectRecords(db, statements.SelectStatement)
	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {

		if err = rows.Scan(&id, &firstName, &email, &username, &instagram, &twitter, &github, &timestamp); err != nil {
			return err
		}

		fmt.Printf("\nChecking username: %s\n", username)
		timeDiff.CalculateDiff(timestamp)

		if instagram {
			if err := runInstagramCheck(email, firstName, username); err != nil {
				return err
			}
		}
		if twitter {
			if err := runTwitterCheck(email, firstName, username); err != nil {
				return err
			}
		}
		if github {
			if err := runGithubCheck(email, firstName, username); err != nil {
				return err
			}
		}
		numUpdated, err := database.UpdateRecords(db, statements.DefaultUpdateStatement, id)
		if err != nil {
			return err
		}
		fmt.Println("Default timestamp update:", numUpdated)
	}

	// Get any error encountered during iteration
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// hanlder gets run every time a Google Cloud CRON Job makes a get request
func handler(w http.ResponseWriter, r *http.Request) {

	keys, ok := r.URL.Query()["auth"]
	if !ok || len(keys[0]) < 1 || keys[0] != os.Getenv("SERVER_PASSWORD") {
		fmt.Fprintf(w, "You are not authorised to access this page")
	} else {
		log.Print("Notify.is: received a request")

		if err := runCheck(); err != nil {
			sentry.CaptureException(err)
			fmt.Println(err)
			fmt.Println("Returning...")
			return
		}
		fmt.Fprintf(w, "Ready to process requests.\n")
	}
}

var db *sql.DB

func init() {

	const (
		port   = 5432
		user   = "postgres"
		dbName = "notify"
	)

	var host = os.Getenv("DB_HOST")
	var password = os.Getenv("DB_PASSWORD")

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, user, password, dbName)

	// Open database connection
	var err error
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		sentry.CaptureException(err)
		fmt.Printf("%v", err)
		fmt.Println("Returning...")
		return
	}
	if err = db.Ping(); err != nil {
		sentry.CaptureException(err)
		log.Fatal(err)
	}
}

func main() {

	// To initialize Sentry's handler, you need to initialize Sentry itself beforehand
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_DSN"),
	}); err != nil {
		fmt.Printf("Sentry initialization failed: %v\n", err)
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	log.Print("Notify.is: starting server...")

	http.HandleFunc("/", sentryHandler.HandleFunc(handler))

	log.Printf("Notify.is: listening on port %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))

}
