package postmark

import (
	"os"

	"github.com/keighl/postmark"
)

var client *postmark.Client

func init() {
	client = postmark.NewClient(os.Getenv("SERVER_TOKEN"), os.Getenv("ACCOUNT_TOKEN"))
}

func SendSuccessEmail(email, name, username, service string) (postmark.EmailResponse, error) {

	var successEmail = postmark.TemplatedEmail{

		From:          "support@notify.is",
		To:            email,
		ReplyTo:       "support@notify.is",
		TemplateId:    22097800,
		TemplateAlias: "success-email",
		TemplateModel: map[string]interface{}{
			"product_url":   "https://notify.is",
			"name":          name,
			"username":      username,
			"service":       service,
			"support_email": "support@notify.is",
			"product_name":  "Notify.is",
		},
	}
	res, err := client.SendTemplatedEmail(successEmail)
	if err != nil {
		return res, err
	}

	return res, nil
}
