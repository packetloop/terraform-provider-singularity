package mesos_singularity

import (
	"github.com/lenfree/go-singularity"
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
	cf := singularity.NewConfig().
		SetHost(c.Host).
		SetPort(c.Port).
		SetRetry(c.Retry).
		Build()

	client := singularity.NewClient(cf)

	return &Conn{
		sclient: client,
	}, nil
}
