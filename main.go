package main

import (
	"log"
	"net"
	"os"
)

var Port = "6969"
var Clients = make(map[string]net.Conn)
var HarumakiServer = os.Getenv("HARUMAKI_SERVER")
var HaruConn net.Conn

func main() {
	// Connect to HarumakiServer (Send Only)
	HaruConn, err := net.Dial("tcp", HarumakiServer)
	if err != nil {
		log.Fatalf("Failed to connect to HarumakiServer: %v", err)
	}
	Clients[HaruConn.RemoteAddr().String()] = HaruConn

	// Start the server
	ln, err := net.Listen("tcp", ":"+Port)
	log.Printf("Listening on port %s", Port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", Port, err)
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Failed to accept connection: %v", err)
		}
		Clients[conn.RemoteAddr().String()] = conn
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	defer func() {
		delete(Clients, conn.RemoteAddr().String())
	}()
	conn.Write([]byte("Welcome to the server\n"))
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("Failed to read from connection: %v", err)
			return
		}

		msg := string(buf[:n])

		log.Printf("message: %s", msg)

		if msg == "exit\r\n" {
			log.Printf("Client requested to close the connection")
			break
		}

		for addr, c := range Clients {
			if addr != conn.RemoteAddr().String() {
				_, err := c.Write([]byte(msg))
				if err != nil {
					log.Printf("Failed to write to connection: %v", err)
				}
			}
		}

		log.Printf("Received %d bytes: %s", n, string(buf[:n]))
	}
	conn.Write([]byte("Goodbye\n"))
}
