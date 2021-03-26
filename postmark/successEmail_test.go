package postmark

import (
	"testing"
)

type addPostmarkTest struct {
	email, name, username, service string
	expected int64
}

var addPostmarkTests = []addPostmarkTest{
	{"support@notify.is", "Support", "username", "Instagram", 0},
}

func TestSendSuccessEmail(t *testing.T) {

	for _, test := range addPostmarkTests{
		if output, _ := SendSuccessEmail(test.email, test.name, test.username, test.service); output.ErrorCode != test.expected {
			t.Errorf("Error code '%d' not equal to expected '%d'", output.ErrorCode, test.expected)
		}
	}

}