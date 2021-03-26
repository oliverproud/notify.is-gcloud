package check

import (
	"fmt"
	"testing"
)

// username is the argument to the function, and the expected stands for the 'result we expect'
type addInstagramTest struct {
	username string
	expected bool
}

var addInstagramTests = []addInstagramTest{
	{"test", false}, // username 'test' is not available, so we expect the check to return false for unavailable
	{"test_available_1234", true}, // username 'test_available_1234098765' is available, so we expect the check to return true for available
}


func TestInstagram(t *testing.T) {

	for _, test := range addInstagramTests{
		output, err := Instagram(test.username)
		if err != nil {
			fmt.Println(warning(err))
		}
		if output != test.expected {
			t.Errorf("Output '%t' not equal to expected '%t'", output, test.expected)
		}
	}

}