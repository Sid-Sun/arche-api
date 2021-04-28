package config

type JWTConfig struct {
	secret string
}

func (j JWTConfig) GetSecret() string {
	return j.secret
}
