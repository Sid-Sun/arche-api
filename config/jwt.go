package config

type JWTConfig struct {
	secret string
	ttl    int
}

func (j JWTConfig) GetSecret() string {
	return j.secret
}

func (j JWTConfig) GetTTL() int {
	return j.ttl
}
