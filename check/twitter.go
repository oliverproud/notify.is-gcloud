package check

import (
	"fmt"
	"github.com/fatih/color"
	"os"
	"strings"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var available bool
// Configure colorized outputs
var warning = color.New(color.FgRed, color.Bold).SprintFunc()
var success = color.New(color.FgGreen, color.Bold).SprintFunc()

// Twitter uses a Go Twitter API to make requests
func Twitter(username string) (bool, error) {
	// Oauth1 required to access Twitter API
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	showParams := &twitter.UserShowParams{ScreenName: username}



	// Get Twitter user
	user, _, err := client.Users.Show(showParams)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "twitter: 215 Bad Authentication data."):
			return false, err
		case strings.Contains(err.Error(), "twitter: 50 User not found."):
			fmt.Printf("Twitter: username %s is %s\n", username, success("AVAILABLE"))
			available = true
		case strings.Contains(err.Error(), "twitter: 63 User has been suspended."):
			fmt.Printf("Twitter: username %s is %s\n", username, warning("SUSPENDED"))
			available = false
		default:
			return available, err
		}
	}

	if user.ID != 0 {
		fmt.Printf("Twitter: username %s is %s \n", username, warning("NOT AVAILABLE"))
		available = false
	}
	return available, nil
}
