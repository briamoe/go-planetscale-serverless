package planetscale

import "encoding/json"

type Config struct {
	Host     string
	Username string
	Password string
}

type Connection struct {
	// Configuration for the connection
	Config *Config

	// Info for the session that the connection uses to queries to PlanetScale.
	// There's no need to parse this as it's passed in it's entirety for each request.
	Session *json.RawMessage
}

// Creates a new connection to PlanetScale
func NewConnection(config *Config) (*Connection, error) {
	c := &Connection{
		Config: config,
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
	r, err := postHTTP[createSessionResponse](c.Config, "CreateSession", struct{}{})
	if err != nil {
		return nil, err
	}

	return r.Session, nil
}
