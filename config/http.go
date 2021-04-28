package config

import "fmt"

type HTTPServerConfig struct {
	host string
	port int
}

func (h HTTPServerConfig) GetListenAddr() string {
	return fmt.Sprintf("%s:%d", h.host, h.port)
}
