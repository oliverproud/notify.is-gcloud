package check

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/go-github/github"
)

var githubAvailable bool

// Github uses a Go Github API to make requests
func Github(username string) (bool, error) {

	client := github.NewClient(nil)

	ctx := context.Background()

	// Get GitHub user
	user, _, err := client.Users.Get(ctx, username)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "404 Not Found"):
			break
		default:
			return githubAvailable, err
		}
	}

	if user != nil {
		githubAvailable = false
		fmt.Printf("GitHub: username %s is %s\n", username, warning("NOT AVAILABLE"))
	} else {
		githubAvailable = true
		fmt.Printf("GitHub: username %s is %s\n", username, success("AVAILABLE"))
	}
	return githubAvailable, nil
}
