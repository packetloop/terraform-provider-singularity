// Package singularity provides simple interface to manage Mesos'
// Singularity jobs and tasks lifecycle.
package singularity

import (
	"strconv"

	"github.com/go-resty/resty"
)

// Client contains Singularity endpoint for http requests
type Client struct {
	Rest *resty.Client
}

// Config contains Singularity HTTP endpoint and configuration for
// retryablehttp client's retry options
type config struct {
	Host  string
	Port  int
	Retry int
}

type ConfigBuilder interface {
	SetPort(int) ConfigBuilder
	SetHost(string) ConfigBuilder
	SetRetry(int) ConfigBuilder
	Build() config
}

// NewConfig returns an empty ConfigBuilder.
func NewConfig() ConfigBuilder {
	return &config{}
}

// SetHost accepts a string and sets the host in config.
func (co *config) SetHost(host string) ConfigBuilder {
	co.Host = host
	return co
}

// SetHost accepts an int and sets the retry count.
func (co *config) SetRetry(r int) ConfigBuilder {
	co.Retry = r
	return co
}

// SetHost accepts an int and sets the port number.
func (co *config) SetPort(port int) ConfigBuilder {
	co.Port = port
	return co
}

// Build method returns a config struct.
func (co *config) Build() config {
	return config{
		Host:  co.Host,
		Port:  co.Port,
		Retry: co.Retry,
	}
}

// NewClient returns Singularity HTTP endpoint.
func NewClient(c config) *Client {
	r := resty.New().
		SetRESTMode().
		SetRetryCount(c.Retry).
		SetHostURL(endpoint(&c))
	return &Client{
		Rest: r,
	}
}

func endpoint(c *config) string {
	// if port is uninitialised, port would be http/80.
	if c.Port == 0 || c.Port == 80 {
		return "http://" + c.Host
	}
	if c.Port == 443 {
		return "https://" + c.Host
	}
	return "http://" + c.Host + ":" + strconv.Itoa(c.Port)
}
