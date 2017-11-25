package singularity

import (
	"testing"
)

func TestNew(t *testing.T) {
	host := "localhost"
	c := Config{
		Host: host,
	}
	client := New(c)
	expected := "http://" + host
	if client.Endpoint != expected {
		t.Errorf("Got %s, expected %s", client.Endpoint, expected)
	}

	port := 8080
	config := Config{
		Host: host,
		Port: port,
	}
	clientConfig := New(config)
	expectedWPort := "http://" + host + ":8080"
	if clientConfig.Endpoint != expectedWPort {
		t.Errorf("Got %s, expected %s", clientConfig.Endpoint, expectedWPort)
	}
}

func TestEndpoint(t *testing.T) {
	config := Config{
		Host: "localhost",
	}
	host := endpoint(&config)
	expected := "http://localhost"
	if host != expected {
		t.Errorf("Got %s, expected %s", host, expected)
	}

	configHTTP := Config{
		Host: "localhost",
		Port: 80,
	}
	hostHTTP := endpoint(&configHTTP)
	expectedHTTP := "http://localhost"
	if hostHTTP != expectedHTTP {
		t.Errorf("Got %s, expected %s", hostHTTP, expectedHTTP)
	}

	configHTTPS := Config{
		Host: "localhost",
		Port: 443,
	}
	hostHTTPS := endpoint(&configHTTPS)
	expectedHTTPS := "https://localhost"
	if hostHTTPS != expectedHTTPS {
		t.Errorf("Got %s, expected %s", hostHTTPS, expectedHTTPS)
	}

}
