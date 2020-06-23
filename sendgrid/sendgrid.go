// Package sendgrid is using SendGrid's Go Library
// https://github.com/sendgrid/sendgrid-go
package sendgrid

import (
	"fmt"
	"log"

	sendgrid "github.com/sendgrid/sendgrid-go"
)

// SendSignup will send the dynamic transactional template when a user signs up
func SendSignup(email, name, username, apiKey, host string) {
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", host)
	request.Method = "POST"
	body := fmt.Sprintf(`{
	 "from":{
			"email":"oliver@oliverproud.com"
	 },
	 "personalizations":[
			{
				 "to":[
						{
							 "email":"%s"
						}
				 ],
				 "dynamic_template_data":{
						"first_name":"%s",
						"username":"%s"

					}
			}
	 ],
	 	"template_id":"d-0449de6d2d8c431a9adb9f079bce3cc7"
	 }`, email, name, username)
	request.Body = []byte(body)

	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		// fmt.Println(response.Body)
		// fmt.Println(response.Headers)
	}
}

// SendSuccess will send a user an email when their username becomes available
func SendSuccess(email, name, username, apiKey, host string) {
	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", host)
	request.Method = "POST"
	body := fmt.Sprintf(`{
	 "from":{
			"email":"oliver@oliverproud.com"
	 },
	 "personalizations":[
			{
				 "to":[
						{
							 "email":"%s"
						}
				 ],
				 "dynamic_template_data":{
						"first_name":"%s",
						"username":"%s"

					}
			}
	 ],
	 	"template_id":"d-8d0bb30d08564ee39fe261040db6f9c3"
	 }`, email, name, username)
	request.Body = []byte(body)

	response, err := sendgrid.API(request)
	if err != nil {
		log.Println(err)
	} else {
		fmt.Println(response.StatusCode)
		// fmt.Println(response.Body)
		// fmt.Println(response.Headers)
	}
}

//
// func main() {
//
// 	name := "Oliver"
// 	email := "owproud@gmail.com"
// 	username := "oliverproud"
//
// 	apiKey := os.Getenv("SENDGRID_API_KEY")
// 	host := "https://api.sendgrid.com"
//
// 	SendSignup(email, name, username, apiKey, host)
// 	SendSuccess(email, name, username, apiKey, host)
// }
