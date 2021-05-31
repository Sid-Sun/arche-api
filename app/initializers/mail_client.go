package initializers

import (
	"fmt"

	"github.com/mailgun/mailgun-go"
	"github.com/sid-sun/arche-api/config"
)

type MailClient interface {
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
	Send(message *mailgun.Message) (mes string, id string, err error)
	Domain() string
}

func InitMGClient(emailConfig *config.EmailConfig) MailClient {
	fmt.Println(emailConfig.GetAPIKey(), emailConfig.GetDomain())
	mgImpl := mailgun.NewMailgun(emailConfig.GetDomain(), emailConfig.GetAPIKey())
	return mgImpl
}
