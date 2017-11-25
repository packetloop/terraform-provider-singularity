// Package singularity provides simple interface to manage Mesos'
// Singularity jobs and tasks lifecycle.
package singularity

import (
	"strconv"

	"github.com/parnurzeal/gorequest"
)

// Client contains Singularity endpoint for http requests
type Client struct {
	Endpoint   string
	SuperAgent gorequest.SuperAgent
}

// Config contains Singularity HTTP endpoint and configuration for
// retryablehttp client's retry options
type Config struct {
	Host  string
	Port  int
	Retry int
}

// New returns Singularity HTTP endpoint.
func New(c Config) *Client {
	a := gorequest.New()
	return &Client{
		Endpoint:   endpoint(&c),
		SuperAgent: *a,
	}
}

func endpoint(c *Config) string {
	// if port is uninitialised, port would be http/80.
	if c.Port == 0 || c.Port == 80 {
		return "http://" + c.Host
	}
	if c.Port == 443 {
		return "https://" + c.Host
	}
	return "http://" + c.Host + ":" + strconv.Itoa(c.Port)
}
