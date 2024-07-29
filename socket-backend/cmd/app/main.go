package main

import "socket-backend/internal/Socket"

func main() {
	const port = "8080"
	server := Socket.NewServer(port)
	server.Run()
}
