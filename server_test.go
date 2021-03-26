package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// name is the name of the test, method is the http request method (e.g. GET, POST),
// password is the server password, wantStatusCode is the response status code we expect
type addServerTest struct {
	name           string
	method         string
	password       string
	wantStatusCode int
}

var addServerTests = []addServerTest{
	{"GET request with password", http.MethodGet, os.Getenv("SERVER_PASSWORD"), http.StatusOK},
	{"GET request without password", http.MethodGet, "", http.StatusUnauthorized},
	{"POST request", http.MethodPost, "", http.StatusMethodNotAllowed},
}

func TestServer(t *testing.T) {

	for _, test := range addServerTests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(test.method, "/?auth=" + test.password, nil)
			responseRecorder := httptest.NewRecorder()

			handler(responseRecorder, request)

			if responseRecorder.Code != test.wantStatusCode {
				t.Errorf("Want status '%d', got '%d'", test.wantStatusCode, responseRecorder.Code)
			}
		})
	}
}
