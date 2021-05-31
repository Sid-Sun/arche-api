package config

import "fmt"

type VerificationEmailConfig struct {
	senderName     string
	senderUsername string
	emailSubject   string
	emailBody      string
	tokenLength    int
}

func newDefaultVEConfig() *VerificationEmailConfig {
	return &VerificationEmailConfig{
		senderName:     "Bouncer",
		senderUsername: "bouncer",
		emailSubject:   "Verify your sign-up and get started!",
		tokenLength:    12,
		emailBody: `
		Hey!

		We noticed that you've signed up to OnlyNotes, we couldn't be more excited to have you join us!
		Please click the link below to activate your account and get started!

		%s

		That's all for now, Cheers!
	`,
	}
}

func (v *VerificationEmailConfig) GetSenderEmail(domain string) string {
	return fmt.Sprintf("%s <%s@%s>", v.senderName, v.senderUsername, domain)
}

func (v *VerificationEmailConfig) GetSubject() string {
	return v.emailSubject
}

func (v *VerificationEmailConfig) GetBody(callbackURL string) string {
	return fmt.Sprintf(v.emailBody, callbackURL)
}

func (v *VerificationEmailConfig) GetTokenLength() int {
	return v.tokenLength
}
