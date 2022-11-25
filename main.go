package main

import (
	"fmt"
	"net"
)

func main() {
	// Start engine here and start the loop
	myGameEngine := NewEngine()

	// Start the tcp lisener
	listener, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		fmt.Println("error listening on TCP port", err)
	}

	defer func() {
		listener.Close()
	}()

	nameGenerator := NewNameGen([]string{RoleExplorer, RoleArchitect, RoleThief, RoleCharlatan})

	fmt.Println("waiting to accept connections...")
	for {
		fmt.Println("waiting for next connection")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error listening on accepting connection", err)
		}
		fmt.Println("got a connection")
		uc, _ := NewUserConnection(conn, myGameEngine)
		player := NewPlayer(uc, nameGenerator.getName())
		myGameEngine.AddPlayer(player)
	}
}
