package client

import "log"

type ReceiveMessageArgs struct {
	ClientId int
	Text     string
}

type ReceiveMessageResponse struct {
}

func (c *Client) ReceiveMessage(req ReceiveMessageArgs, res *ReceiveMessageResponse) error {
	if req.ClientId != c.Id {
		log.Printf("[%d] : %s \n", req.ClientId, req.Text)
	}
	return nil
}
