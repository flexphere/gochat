package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

const ENDSEQ = "exit\r\n"

type Client struct {
	conn net.Conn
	emit EventListener
	Name string
}

func NewClient(conn net.Conn, el EventListener) *Client {
	return &Client{
		conn: conn,
		emit: el,
	}
}

func (c *Client) Listen() {
	defer func() { c.emit.Disconnect <- c.Address() }()

	c.emit.Connect <- c
	for {
		buf := make([]byte, 1024)
		n, err := c.conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		msg := string(buf[:n])
		if msg == ENDSEQ {
			return
		}

		if c.Name == "" {
			c.Name = strings.TrimSpace(msg)
			continue
		}

		c.emit.Broadcast <- &BroadcastMessage{
			From:    c.Address(),
			Message: fmt.Sprintf("%s: %s", c.Name, msg),
		}
	}
}

func (c *Client) Write(msg string) {
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		log.Printf("Failed to write to connection: %v", err)
	}
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Address() string {
	return c.conn.RemoteAddr().String()
}
