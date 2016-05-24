package client

type Client struct {
	server string
	nodeID string

	clientOption  *ClientOption
	profileOption *ProfileOption
}

func NewClient(server, nodeID string, clientOption *ClientOption) *Client {
	if clientOption == nil {
		clientOption = NewClientOption()
	}
	c := &Client{
		server:       server,
		nodeID:       nodeID,
		clientOption: clientOption,
	}
	return c
}

func (c *Client) Run() {
	go c.run()
}

func (c *Client) run() {
	// connect to server

	// run profile task

	// post task file to server

	// remove tmp file
}
