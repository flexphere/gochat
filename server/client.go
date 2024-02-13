package server

import (
	"log"
	"net"
)

type Client struct {
	Conn net.Conn
	Emit EventListener
}

func NewClient(conn net.Conn, el EventListener) *Client {
	return &Client{
		Conn: conn,
		Emit: el,
	}
}

func (c *Client) Listen() {
	defer func() { c.Emit.Disconnect <- c.Address() }()

	c.Emit.Connect <- c
	for {
		buf := make([]byte, 1024)
		n, err := c.Conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		msg := string(buf[:n])
		if msg == "exit\r\n" {
			return
		}

		c.Emit.Broadcast <- BroadcastMessage{
			From:    c.Address(),
			Message: msg,
		}
	}
}

func (c *Client) Write(msg string) {
	_, err := c.Conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Failed to write to connection: %v", err)
	}
}

func (c *Client) Close() {
	c.Conn.Close()
}

func (c *Client) Address() string {
	return c.Conn.RemoteAddr().String()
}
