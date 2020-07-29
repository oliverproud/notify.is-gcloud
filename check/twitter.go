package check

import (
	"fmt"
	"os"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

var available bool

// Twitter uses a Go Twitter API to make requests
func Twitter(username string) (bool, error) {
	config := oauth1.NewConfig(os.Getenv("CONSUMER_KEY"), os.Getenv("CONSUMER_SECRET"))
	token := oauth1.NewToken(os.Getenv("TOKEN"), os.Getenv("TOKEN_SECRET"))
	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	showParams := &twitter.UserShowParams{ScreenName: username}

	user, _, err := client.Users.Show(showParams)
	if err != nil {
		if err.Error() == "twitter: 50 User not found." {
			available = true
			fmt.Printf("Twitter: username %s available\n", username)
		} else {
			return available, err
		}
	}

	if user.ID != 0 {
		fmt.Printf("Twitter: username %s NOT available: ID %d\n", username, user.ID)
		available = false
	}
	return available, nil
}
