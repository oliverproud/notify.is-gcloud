package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"notify.is-go/check"
	"notify.is-go/postmark"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User database struct
type User struct {
	gorm.Model
	Instagram, Twitter, Github bool
	FirstName, Email, Username string
	ID                         string    `gorm:"default:uuid_generate_v4()"`
	Timestamp                  time.Time `gorm:"default:timezone('utc'::text, now())"`
}

// Checks Instagram, sends email, updates database
func runInstagramCheck(email, firstName, username string, user User) error {
	instagramAvailable, err := check.Instagram(username)
	if err != nil {
		return err
	}

	if instagramAvailable {
		resp, err := postmark.SendSuccessEmail(email, firstName, username, "Instagram")
		if err != nil {
			return err
		}

		result := db.Model(&user).Updates(&User{Instagram: false, Timestamp: time.Now()})
		fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)
		fmt.Println("Number of records updated:", result.RowsAffected)
	}
	return nil
}

// Checks Twitter, sends email, updates database
func runTwitterCheck(email, firstName, username string, user User) error {
	twitterAvailable, err := check.Twitter(username)
	if err != nil {
		return err
	}

	if twitterAvailable {
		resp, err := postmark.SendSuccessEmail(email, firstName, username, "Twitter")
		if err != nil {
			return err
		}

		result := db.Model(&user).Updates(&User{Twitter: false, Timestamp: time.Now()})
		fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)
		fmt.Println("Number of records updated:", result.RowsAffected)
	}
	return nil
}

// Checks GitHub, sends email, updates database
func runGithubCheck(email, firstName, username string, user User) error {

	githubAvailable, err := check.Github(username)
	if err != nil {
		return err
	}

	if githubAvailable {
		resp, err := postmark.SendSuccessEmail(email, firstName, username, "GitHub")
		if err != nil {
			return err
		}

		result := db.Model(&user).Updates(&User{Github: false, Timestamp: time.Now()})
		fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)
		fmt.Println("Number of records updated:", result.RowsAffected)
	}
	return nil
}

// Selects records from database runs checks on them
func runCheck() error {
	log.Println("Starting check...")

	rows, err := db.Model(&User{}).Where("EXTRACT(EPOCH FROM ((now() at time zone 'utc') - timestamp)) > 43200").Rows()
	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var user User
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		db.ScanRows(rows, &user)

		fmt.Printf("\nChecking username: %s\n", user.Username)
		fmt.Printf("Time since last check: %v\n", time.Since(user.Timestamp))

		if user.Instagram {
			if err := runInstagramCheck(user.Email, user.FirstName, user.Username, user); err != nil {
				return err
			}
		}
		if user.Twitter {
			if err := runTwitterCheck(user.Email, user.FirstName, user.Username, user); err != nil {
				return err
			}
		}
		if user.Github {
			if err := runGithubCheck(user.Email, user.FirstName, user.Username, user); err != nil {
				return err
			}
		}

		// Default timestamp update
		result := db.Model(&user).Updates(&User{Timestamp: time.Now()})
		fmt.Println("Default timestamp update:", result.RowsAffected)
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

var db *gorm.DB

func init() {

	const (
		port   = 5432
		dbUser = "postgres"
		dbName = "notify"
	)

	var host = os.Getenv("DB_HOST")
	var password = os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require", host, port, dbUser, password, dbName)

	// Open database connection
	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
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
