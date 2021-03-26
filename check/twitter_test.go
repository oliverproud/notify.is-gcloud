package check

import (
	"fmt"
	"testing"
)

// username is the argument to the function, and the expected stands for the 'result we expect'
type addTwitterTest struct {
	username string
	expected bool
}

var addTwitterTests = []addTwitterTest{
	{"oliverwproud", false},
	{"oliverproud_available", true},
	{"suspendedaccounthere", false},
}


func TestTwitter(t *testing.T) {

	for _, test := range addTwitterTests{
		output, err := Twitter(test.username)
		if err != nil {
			fmt.Println(warning(err))
		}
		if output != test.expected {
			t.Errorf("Output '%t' not equal to expected '%t'", output, test.expected)
		}
	}
}