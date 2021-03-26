package check

import (
	"context"
	"fmt"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
	"log"
	"strings"
	"time"
)

// instagramAvailable let's the main package know if a username is available
var instagramAvailable bool
var nodes []*cdp.Node

const (
	urlStr           = `https://www.instagram.com/accounts/emailsignup/`
	usernameSelector = `//input[@name="username"]`
	bodySelector     = `/html/body`
	spriteSelector   = `//span[contains(@class,'gBp1f')]`
)

// Instagram runs the headless browser that checks Instagram
func Instagram(username string) (bool, error) {

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.90 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	// also set up a custom logger
	ctx, cancel := chromedp.NewContext(allocCtx, /*chromedp.WithDebugf(log.Printf),*/ chromedp.WithLogf(log.Printf))
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
				fmt.Printf("Instagram: username %s is %s\n", username, success("AVAILABLE"))
			} else if strings.Contains(nodes[0].Attributes[i], "coreSpriteInputError") {
				instagramAvailable = false
				fmt.Printf("Instagram: username %s is %s\n", username, warning("NOT AVAILABLE"))
			}
		}
	} else {
		return instagramAvailable, fmt.Errorf("no sprites returned")
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