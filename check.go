// Command submit is a chromedp example demonstrating how to fill out and
// submit a form.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"username-check/sendgrid"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

// XHRResponse handles the XHR JSON data coming in from Instagram
type XHRResponse struct {
	AccountCreated bool `json:"account_created"`
	Errors         struct {
		Username []struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"username"`
	} `json:"errors"`
}

var parseXHR XHRResponse
var available bool

func main() {

	// SendGrid API key
	apiKey := os.Getenv("SENDGRID_API_KEY")
	host := "https://api.sendgrid.com"

	firstName := flag.String("name", "Oliver", "users first name")
	email := flag.String("email", "owproud@gmail.com", "users email address")
	username := flag.String("username", "oliverproud", "users desired username")

	flag.Parse()

	fmt.Println(*firstName)
	fmt.Println(*email)
	fmt.Println(*username)

	run := true
	if run {

		// create context
		ctx, cancel := chromedp.NewContext(context.Background())
		defer cancel()

		// run task list
		err := chromedp.Run(ctx, submit(ctx, `https://www.instagram.com/accounts/emailsignup/`, `//input[@name="username"]`, *email, *firstName, *username, apiKey, host))
		if err != nil {
			log.Fatal(err)
		}

	}
}

func submit(ctx context.Context, urlstr, selector, email, name, username, apiKey, host string) chromedp.Tasks {

	chromedp.ListenTarget(ctx, func(event interface{}) {
		if event, ok := event.(*network.EventResponseReceived); ok {

			if event.Type != "XHR" {
				return
			}

			go func() {
				// print response body
				c := chromedp.FromContext(ctx)
				rbp := network.GetResponseBody(event.RequestID)
				body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Println(err)
				}

				// Check XHR response body for correct data
				if strings.HasPrefix(string(body), `{"account_created"`) {

					// Parse JSON data
					json.Unmarshal([]byte(body), &parseXHR)
					if parseXHR.Errors.Username != nil {
						fmt.Println("Username is taken")
					} else {
						fmt.Println("Username is available")
						available = true
					}

					if available == true {
						sendgrid.SendSuccess(email, name, username, apiKey, host)
					}
				}
			}()
		}
	})

	return chromedp.Tasks{

		network.Enable(),
		chromedp.Navigate(urlstr),
		chromedp.WaitVisible(selector),
		chromedp.SendKeys(selector, username),
		chromedp.Click(`//*[@id="react-root"]/section/main/div/article/div/div[1]/div/form/div[7]/div/button`, chromedp.BySearch),
		chromedp.Sleep(time.Second * 1),
	}
}
