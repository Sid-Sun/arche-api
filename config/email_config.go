package config

type EmailConfig struct {
	domain string
	apiKey string
}

func (e *EmailConfig) GetDomain() string {
	return e.domain
}

func (e *EmailConfig) GetAPIKey() string {
	return e.apiKey
}
