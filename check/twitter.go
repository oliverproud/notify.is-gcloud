package check

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var username string

func init() {
	***REMOVED***("CONSUMER_KEY", "***REMOVED***")
	***REMOVED***("CONSUMER_SECRET", "***REMOVED***")
	***REMOVED***("TOKEN", "***REMOVED***")
	***REMOVED***("TOKEN_SECRET", "***REMOVED***")

	username = "oliverwproud"
}

// TwitterAPI uses a Golang Twitter API to make requests
func TwitterAPI() {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	showParams := &twitter.UserShowParams{ScreenName: username}

	user, _, err := client.Users.Show(showParams)
	if err != nil {
		if err.Error() == "twitter: 50 User not found." {
			// Do database stuff in here
			fmt.Println(err)
		} else {
			fmt.Println(err)
			return
		}
	}

	if user.ID != 0 {
		// This is where success email and database deletion would happen
		fmt.Println("USER:", user.ID)
	}

}
