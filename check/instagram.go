package check

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// instagramAvailable let's the main package know if a username is available
var instagramAvailable bool
var nodes []*cdp.Node

const (
	urlStr           = `https://www.instagram.com/accounts/emailsignup/`
	usernameSelector = `//input[@name="username"]`
	bodySelector     = `//div[contains(@class,'tbpKJ')]`
	spriteSelector   = `//span[contains(@class,'gBp1f')]`
)

// Instagram runs the headless browser that checks Instagram
func Instagram(username string) (bool, error) {

	// create chrome instance
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	task, err := submit(urlStr, usernameSelector, bodySelector, spriteSelector, username)
	if err != nil {
		return instagramAvailable, err
	}

	// run task
	if err = chromedp.Run(ctx, task); err != nil {
		return instagramAvailable, err
	}

	if len(nodes) != 0 {
		for i := range nodes[0].Attributes {
			if strings.Contains(nodes[0].Attributes[i], "coreSpriteInputAccepted") {
				instagramAvailable = true
				fmt.Printf("Instagram: username %s is available\n", username)
			} else if strings.Contains(nodes[0].Attributes[i], "coreSpriteInputError") {
				instagramAvailable = false
				fmt.Printf("Instagram: username %s is NOT available\n", username)
			}
		}
	} else {
		return instagramAvailable, fmt.Errorf("No sprites returned")
	}

	return instagramAvailable, nil
}

func submit(urlStr, usernameSelector, bodySelector, spriteSelector, username string) (chromedp.Tasks, error) {

	return chromedp.Tasks{
		chromedp.Navigate(urlStr),
		chromedp.WaitVisible(usernameSelector),
		chromedp.SendKeys(usernameSelector, username),
		// XPath BySearch just looking for Body
		chromedp.WaitVisible(bodySelector),
		chromedp.Click(bodySelector),
		// XPath with 'contains' function looking for class
		chromedp.WaitVisible(spriteSelector),
		chromedp.Nodes(spriteSelector, &nodes, chromedp.AtLeast(0)),
	}, nil
}
