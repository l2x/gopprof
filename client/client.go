package client

// Client is a gopprof client
type Client struct {
	server string
	nodeID string
}

// NewClient return client
func NewClient(server, nodeID string) *Client {
	c := &Client{
		server: server,
		nodeID: nodeID,
	}
	return c
}

// Run an client
func (c *Client) Run() error {
	if err := c.register(); err != nil {
		return err
	}
	go c.run()
	return nil
}

func (c *Client) register() error {
	return nil
}

func (c *Client) run() {
	// connect to server

	// run profile task

	// post task file to server

	// remove tmp file
}
