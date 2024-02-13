package main

import "chatapp/server"

var Port = "6969"

func main() {
	server := server.NewServer(Port)
	server.Start()
}
