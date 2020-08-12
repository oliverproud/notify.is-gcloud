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

// Instagram runs the headless browser that checks Instagram
func Instagram(username string) (bool, error) {

	// create chrome instance
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithLogf(log.Printf))
	defer cancel()

	// create a timeout
	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	task, err := submit(`https://www.instagram.com/accounts/emailsignup/`, `//input[@name="username"]`, username)
	if err != nil {
		return instagramAvailable, err
	}

	// run task
	if err = chromedp.Run(ctx, task); err != nil {
		return instagramAvailable, err
	}

	if len(nodes) != 0 {
		// fmt.Printf("RESULT: %v\n", nodes[0])
		// fmt.Println("Sprite Type:", nodes[0].Attributes)
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

func submit(urlstr, selector, username string) (chromedp.Tasks, error) {

	return chromedp.Tasks{
		chromedp.Navigate(urlstr),
		chromedp.Sleep(time.Second),
		// chromedp.WaitVisible(selector),
		chromedp.SendKeys(selector, username),
		chromedp.Sleep(time.Second),
		// chromedp.WaitVisible(`//*[@id="react-root"]/section/main/div`),
		// chromedp.Click(`//*[@id="react-root"]/section/main/div/article/div/div[1]/div/form/div[7]/div/button`, chromedp.BySearch),
		chromedp.Click(`//*[@id="react-root"]/section/main/div`, chromedp.BySearch),
		chromedp.Sleep(time.Second),
		// chromedp.WaitVisible(`//*[@id="react-root"]/section/main/div/article/div/div[1]/div/form/div[5]/div/div/span`),
		chromedp.Nodes(`//*[@id="react-root"]/section/main/div/article/div/div[1]/div/form/div[5]/div/div/span`, &nodes, chromedp.AtLeast(0)),
	}, nil
}
