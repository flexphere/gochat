package server

import (
	"log"
	"net"
)

type Server struct {
	Port          string
	Clients       map[string]*Client
	eventListener EventListener
}

type EventListener struct {
	Connect    chan *Client
	Disconnect chan string
	Broadcast  chan *BroadcastMessage
}

type BroadcastMessage struct {
	From    string
	Message string
}

func NewServer(port string) *Server {
	return &Server{
		Port:    port,
		Clients: make(map[string]*Client),
		eventListener: EventListener{
			Connect:    make(chan *Client),
			Disconnect: make(chan string),
			Broadcast:  make(chan *BroadcastMessage),
		},
	}
}

func (s *Server) Start() {
	ln, err := net.Listen("tcp", ":"+s.Port)
	log.Printf("Listening on port %s", s.Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", s.Port, err)
	}

	go s.eventHandler()
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
		}

		go NewClient(conn, s.eventListener).Listen()
	}
}

func (s *Server) eventHandler() {
	for {
		select {
		case client := <-s.eventListener.Connect:
			client.Write("Welcome to the server\n\nusername:")
			s.Clients[client.Address()] = client

		case addr := <-s.eventListener.Disconnect:
			s.Clients[addr].Write("Goodbye\n")
			s.Clients[addr].Close()
			delete(s.Clients, addr)
			log.Printf("Closed connection from %s", addr)

		case data := <-s.eventListener.Broadcast:
			log.Printf("message from (%s): %s", data.From, data.Message)
			for addr, c := range s.Clients {
				if addr == data.From {
					continue
				}
				c.Write(data.Message)
			}
		}
	}
}
