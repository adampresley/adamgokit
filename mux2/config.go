package mux2

type MuxConfig interface {
	GetHost() string
}

type Config struct {
	Host string `flag:"host" env:"HOST" default:"localhost:8081" description:"The address and port to bind the HTTP server to"`
}

func (c Config) GetHost() string {
	return c.Host
}
