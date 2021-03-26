package main

import (
	"fmt"
	"github.com/fatih/color"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"log"
	"math/rand"
	"net/http"
	"notify.is-go/postmark"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"notify.is-go/check"
)

// User database struct
type User struct {
	gorm.Model
	Instagram, Twitter, Github bool
	FirstName, Email, Username string
	ID                         string    `gorm:"default:uuid_generate_v4()"`
}

// Configure colorized outputs
var warning = color.New(color.FgRed, color.Bold).SprintFunc()

// Checks Instagram, sends email, updates database
func runInstagramCheck(email, firstName, username string, user User) error {
	instagramAvailable, err := check.Instagram(username)
	if err != nil {
		return err
	}

	if instagramAvailable {

		// returning any error will rollback changes
		err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Model(&user).Updates(map[string]interface{}{"instagram": false}).Error; err != nil {
				return err
			}

			resp, err := postmark.SendSuccessEmail(email, firstName, username, "Instagram")
			if err != nil {
				return err
			}
			if resp.ErrorCode != 0 {
				return fmt.Errorf("Postmark error: %v %s\n", resp.ErrorCode, resp.Message)
			}
			fmt.Println("Number of records updated:", 1)
			fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)

			// returning nil will commit the whole transaction
			return nil
		})
		if err != nil {
			return err
		}
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

		err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Model(&user).Updates(map[string]interface{}{"twitter": false}).Error; err != nil {
				return err
			}

			resp, err := postmark.SendSuccessEmail(email, firstName, username, "Twitter")
			if err != nil {
				return err
			}
			if resp.ErrorCode != 0 {
				return fmt.Errorf("Postmark error: %v %s\n", resp.ErrorCode, resp.Message)
			}
			fmt.Println("Number of records updated:", 1)
			fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)

			return nil
		})
		if err != nil {
			return err
		}
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

		err := db.Transaction(func(tx *gorm.DB) error {

			if err := tx.Model(&user).Updates(map[string]interface{}{"github": false}).Error; err != nil {
				return err
			}

			resp, err := postmark.SendSuccessEmail(email, firstName, username, "GitHub")
			if err != nil {
				return err
			}
			if resp.ErrorCode != 0 {
				return fmt.Errorf("Postmark error: %v %s\n", resp.ErrorCode, resp.Message)
			}
			fmt.Println("Number of records updated:", 1)
			fmt.Printf("Postmark response: %v %s\n", resp.ErrorCode, resp.Message)

			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// Selects records from database runs checks on them
func runCheck() error {
	log.Println("Starting check...")

	// 43200 seconds == 12 hours
	// 86400 seconds == 24 hours
	rows, err := db.Model(&User{}).Where("EXTRACT(EPOCH FROM ((now() at time zone 'utc') - updated_at)) > 86400").Rows()
	if err != nil {
		return err
	}

	for rows.Next() {
		var user User
		// ScanRows is a method of `gorm.DB`, it can be used to scan a row into a struct
		err := db.ScanRows(rows, &user)
		if err != nil {
			return err
		}

		fmt.Printf("\nChecking username: %s\n", user.Username)
		fmt.Printf("Time since last check: %v\n", time.Since(user.UpdatedAt))

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

		// Update time user was last checked
		result := db.Model(&user).Updates(&User{})
		fmt.Println("Default timestamp update. Rows affected:", result.RowsAffected)

		if user.Instagram {
			// sleep after each Instagram check
			sleepTime := rand.Intn(60 - 10) + 10
			fmt.Printf("Sleeping for %d seconds after Instagram check\n", sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Second)
		}
	}

	// Get any error encountered during iteration
	if err = rows.Err(); err != nil {
		return err
	}
	return nil
}

// handler gets run every time a Google Cloud CRON Job makes a get request
func handler(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		keys, ok := r.URL.Query()["auth"]
		if !ok || len(keys[0]) < 1 || keys[0] != os.Getenv("SERVER_PASSWORD") {
			http.Error(w, "You are not authorised to access this page.", http.StatusUnauthorized)
		} else {
			log.Print("Notify.is: received a request")

			if err := runCheck(); err != nil {
				sentry.CaptureException(err)
				fmt.Println(warning(err))
				fmt.Println("Returning...")
				return
			}
			_, err := fmt.Fprintf(w, "Ready to process requests.\n")
			if err != nil {
				return
			}
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

var db *gorm.DB

func init() {

	var port = 5432
	var dbUser = os.Getenv("DB_USER")
	var dbName = os.Getenv("DB_NAME")
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
		fmt.Printf(warning("Sentry initialization failed: ") + "%v\n", err)
	}

	// Create an instance of sentryhttp
	sentryHandler := sentryhttp.New(sentryhttp.Options{})

	http.HandleFunc("/", sentryHandler.HandleFunc(handler))

	log.Print("Notify.is: starting server...")

	log.Printf("Notify.is: listening on port: %s", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", os.Getenv("PORT")), nil))
}
