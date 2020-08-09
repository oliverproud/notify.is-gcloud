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
)

// InstagramResponse handles the XHR JSON data coming in from Instagram
type InstagramResponse struct {
	AccountCreated bool `json:"account_created"`
	Errors         struct {
		Username []struct {
			Message string `json:"message"`
			Code    string `json:"code"`
		} `json:"username"`
	} `json:"errors"`
}

// InstagramAvailable let's the main package know if a username is available
var InstagramAvailable bool
var parseInstagramResponse InstagramResponse

// Instagram runs the headless browser that checks Instagram
func Instagram(username string) error {

	// create context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	task, err := submit(ctx, `https://www.instagram.com/accounts/emailsignup/`, `//input[@name="username"]`, username)
	if err != nil {
		return err
	}
	// run task
	if err = chromedp.Run(ctx, task); err != nil {
		return err
	}
	return nil
}

func submit(ctx context.Context, urlstr, selector, username string) (chromedp.Tasks, error) {

	chromedp.ListenTarget(ctx, func(event interface{}) {
		if event, ok := event.(*network.EventResponseReceived); ok {

			// Ignore non XHR responses
			if event.Type != "XHR" {
				return
			}

			go func() {
				// Get response body
				c := chromedp.FromContext(ctx)
				rbp := network.GetResponseBody(event.RequestID)
				body, err := rbp.Do(cdp.WithExecutor(ctx, c.Target))
				if err != nil {
					fmt.Println(err)
					return
				}

				// Check XHR response body for correct data
				if strings.HasPrefix(string(body), `{"account_created"`) {

					// Parse JSON data
					json.Unmarshal([]byte(body), &parseInstagramResponse)
					if parseInstagramResponse.Errors.Username != nil {
						fmt.Printf("Instagram: username %s is NOT available\n", username)
						InstagramAvailable = false
					} else {
						fmt.Printf("Instagram: username %s is available\n", username)
						InstagramAvailable = true
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
