package planetscale

import (
	"encoding/json"
	"net/http"
)

type Config struct {
	Host     string
	Username string
	Password string

	Transport http.RoundTripper
}

type Connection struct {
	// Configuration for the connection
	Config *Config

	// Info for the session that the connection uses to queries to PlanetScale.
	// There's no need to parse this as it's passed in it's entirety for each request.
	Session *json.RawMessage

	client *http.Client
}

// Creates a new connection to PlanetScale
func NewConnection(config *Config) (*Connection, error) {
	t := http.DefaultTransport
	if config.Transport != nil {
		t = config.Transport
	}

	c := &Connection{
		Config: config,
		client: &http.Client{
			Transport: t,
		},
	}

	s, err := c.createSession()
	if err != nil {
		return nil, err
	}
	c.Session = s

	return c, nil
}

type createSessionResponse struct {
	Session *json.RawMessage `json:"session"`
}

func (c *Connection) createSession() (*json.RawMessage, error) {
	// creates a session to be reused across requests
	r, err := postHTTP[createSessionResponse](c, "CreateSession", struct{}{})
	if err != nil {
		return nil, err
	}

	return r.Session, nil
}
