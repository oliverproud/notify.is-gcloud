package check

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"notify.is-go/sendgrid"
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

// Available let's the main package know if a username is available
var Available bool
var parseXHR XHRResponse

// RunHeadless runs the headless browser than checks Instagram
func RunHeadless(email, name, username string) error {

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	task, err := submit(ctx, `https://www.instagram.com/accounts/emailsignup/`, `//input[@name="username"]`, email, name, username)
	if err != nil {
		return err
	}
	// run task list
	err = chromedp.Run(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func submit(ctx context.Context, urlstr, selector, email, name, username string) (chromedp.Tasks, error) {

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
					return
				}

				// Check XHR response body for correct data
				if strings.HasPrefix(string(body), `{"account_created"`) {

					// Testing
					//fmt.Println("BODY:", string(body))

					// Parse JSON data
					json.Unmarshal([]byte(body), &parseXHR)
					if parseXHR.Errors.Username != nil {
						fmt.Printf("Username: %s is NOT available\n", username)
						Available = false
					} else {
						fmt.Printf("Username: %s is available\n", username)
						Available = true
						sendgrid.SendEmail(email, name, username, "", "success")
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
		chromedp.Sleep(time.Second * 2),
	}, nil
}
