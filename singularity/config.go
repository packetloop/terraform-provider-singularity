package singularity

import (
	"github.com/lenfree/go-mesos-singularity"
)

// Conn is the client connection manager for singularity provider.
// It holds the connection information such as API endpoint to interface with.
type Conn struct {
	sclient *singularity.Client
}

// Config holds the provider configuration, and delivers a populated
// singularity connection based off the contained settings.
type Config struct {
	Host  string
	Port  int
	Retry int
}

// Client returns a new client for accessing Singularity Rest API.
// We don't do any authorisation as of the moment. Hence, this block
// is simple.
func (c *Config) Client() (*Conn, error) {
	config := singularity.Config{
		Host:  c.Host,
		Port:  c.Port,
		Retry: c.Retry,
	}
	client := singularity.New(config)

	return &Conn{
		sclient: client,
	}, nil
}
