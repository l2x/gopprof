package client

type Client struct {
	server string
	nodeID string

	option *Option
}

func NewClient(server, nodeID string, option *Option) *Client {
	if option == nil {
		option = NewOption()
	}
	c := &Client{
		server: server,
		nodeID: nodeID,
		option: option,
	}
	return c
}

func (c *Client) Run() {
	go c.run()
}

func (c *Client) run() {
}
