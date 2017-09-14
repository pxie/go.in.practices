package lib

import (
	"log"
)

// Client interface
type Client interface {
	IsConnected() bool
	Connect() error
	Disconnect(quiesce uint)
	GetID() string
}

// client implements the Client interface
type client struct {
	lastSent int64
	ID       string
	opts     []string
}

// NewClient will create one Client instance
func NewClient() Client {
	c := &client{}
	c.lastSent = 5
	c.ID = "abc"
	c.opts = []string{"123", "xyz"}

	return c
}

func (c *client) IsConnected() bool {
	log.Println("is connected")
	return true
}

func (c *client) Connect() error {
	log.Println("connect")
	return nil
}

func (c *client) Disconnect(quiesce uint) {
	log.Println("disconnect")
}

func (c *client) GetID() string {
	log.Println("get id:", c.ID)
	return c.ID
}
