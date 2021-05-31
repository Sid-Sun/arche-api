package initializers

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/sid-sun/arche-api/config"
)

type MailClient interface {
	NewMessage(from, subject, text string, to ...string) *mailgun.Message
	Send(ctx context.Context, message *mailgun.Message) (mes string, id string, err error)
	Domain() string
}

func InitMGClient(emailConfig *config.EmailConfig) MailClient {
	mgImpl := mailgun.NewMailgun(emailConfig.GetDomain(), emailConfig.GetAPIKey())
	mgImpl.SetAPIBase(mailgun.APIBaseEU)
	return mgImpl
}
