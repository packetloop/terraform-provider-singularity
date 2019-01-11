package singularity

import (
	"testing"
)

func TestBuild(t *testing.T) {

	port := 80
	retry := 2
	host := "localhost"

	expected := config{
		Host:  "localhost",
		Port:  port,
		Retry: retry,
	}
	c := NewConfig().
		SetHost(host).
		SetPort(port).
		SetRetry(retry).
		Build()
	if c != expected {
		t.Errorf("Got %v, expected %v", c, expected)
	}
}

func TestSetPort(t *testing.T) {
	port := 80

	expected := config{
		Port: port,
	}

	c := config{
		Port: port,
	}

	r := NewConfig().
		SetPort(port).
		Build()

	if r != expected {
		t.Errorf("Got %v, expected %v", c, expected)
	}
}

func TestSetHost(t *testing.T) {
	host := "localhost"

	expected := config{
		Host: host,
	}

	c := config{
		Host: host,
	}

	r := NewConfig().
		SetHost(host).
		Build()

	if r != expected {
		t.Errorf("Got %v, expected %v", c, expected)
	}
}

func TestSetRetry(t *testing.T) {
	retry := 2

	expected := config{
		Retry: retry,
	}

	c := config{
		Retry: retry,
	}

	r := NewConfig().
		SetRetry(retry).
		Build()

	if r != expected {
		t.Errorf("Got %v, expected %v", c, expected)
	}
}

func TestNewClient(t *testing.T) {

	c := NewConfig().SetHost("localhost").SetPort(80).SetRetry(2).Build()

	expected := "http://localhost"
	client := NewClient(c)
	if client.Rest.HostURL != expected {
		t.Errorf("Got %v, expected %v", c, expected)
	}
}

func TestEndpoint(t *testing.T) {
	c := config{
		Host: "localhost",
	}
	host := endpoint(&c)
	expected := "http://localhost"
	if host != expected {
		t.Errorf("Got %s, expected %s", host, expected)
	}

	cHTTP := config{
		Host: "localhost",
		Port: 80,
	}
	hostHTTP := endpoint(&cHTTP)
	expectedHTTP := "http://localhost"
	if hostHTTP != expectedHTTP {
		t.Errorf("Got %s, expected %s", hostHTTP, expectedHTTP)
	}

	cHTTPS := config{
		Host: "localhost",
		Port: 443,
	}
	hostHTTPS := endpoint(&cHTTPS)
	expectedHTTPS := "https://localhost"
	if hostHTTPS != expectedHTTPS {
		t.Errorf("Got %s, expected %s", hostHTTPS, expectedHTTPS)
	}

}
