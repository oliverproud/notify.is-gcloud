package check

import (
	"fmt"
	"testing"
)

// username is the argument to the function, and the expected stands for the 'result we expect'
type addGithubTest struct {
	username string
	expected bool
}

var addGithubTests = []addGithubTest{
	{"oliverproud", false},
	{"oliverproud-available", true},
}


func TestGithub(t *testing.T) {

	for _, test := range addGithubTests{
		output, err := Github(test.username)
		if err != nil {
			fmt.Println(warning(err))
		}
		if output != test.expected {
			t.Errorf("Output '%t' not equal to expected '%t'", output, test.expected)
		}
	}

}