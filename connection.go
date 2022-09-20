package planetscale

type Config struct {
	Host     string
	Username string
	Password string
}

type session struct{}

type Connection struct {
	// Configuration for the connections
	Config *Config

	// Info for the session that the connection uses to make calls to the database.
	Session *session
}

func NewConnection(config *Config) (*Connection, error) {
	c := &Connection{
		Config: config,
	}

	s, err := c.createSession()
	if err != nil {
		// TODO: handle
		return nil, err
	}
	c.Session = s

	return c, nil
}

type createSessionResponse struct {
	Session *session `json:"session"`
}

func (c *Connection) createSession() (*session, error) {
	r, err := postHTTP[createSessionResponse](c.Config, "CreateSession", struct{}{})
	if err != nil {
		return nil, err
	}

	return r.Session, nil
}
