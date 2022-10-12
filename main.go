package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		fmt.Println("error listening on TCP port", err)
	}

	defer func() {
		listener.Close()
	}()

	fmt.Println("waiting to accept connections...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error listening on accepting connection", err)
		}
		userConn := UserConnection{conn: conn}
		if err := handleConnection(&userConn); err != nil {
			fmt.Println("error handling connection", err)
		}

		conn.Close()
	}
}

func handleConnection(conn *UserConnection) error {
	var err error
	conn.Send(IntroString())
	conn.Prompt()
	r := bufio.NewReader(conn.conn)
	data, err := r.ReadString('\n')
	if err != nil {
		return err
	}

	for {
		out, _ := handleInput(data)
		if err := conn.Send(out); err == nil {
			return err
		}
	}
	return nil
}

func handleInput(in string) (string, error) {
	cleanIn := strings.TrimSpace(in)
	switch strings.ToLower(cleanIn) {
	case "l", "look":
		return "You take a look around your current room. It is empty. You should quit.\n", nil
	case "i", "inspect":
		return "You must specify what you want to inspect.\n", nil
	case "u", "use":
		return "You must specify what you want to use.\n", nil
	case "q", "quit":
		return "Quitter!!!\n", nil
	default:
		return fmt.Sprintf("I don't know what %s that is.\n", in), fmt.Errorf("got command I cannot handle")
	}
}

// On user input, the command should
// 1. Optionally update the state of the world
// 1. Queue up messages to send to the user
// 1. Print a message to the user. Note that ther emust always be anew line

func IntroString() string {
	return strings.Join([]string{
		`
///////////////////////////////////////////////////////////////////////////////////////////
//  _________          _______    _______  _______  _______  _______  _______  _______ 
//  \__   __/|\     /|(  ____ \  (  ____ \(  ____ \(  ____ \(  ___  )(  ____ )(  ____ \
//     ) (   | )   ( || (    \/  | (    \/| (    \/| (    \/| (   ) || (    )|| (    \/
//     | |   | (___) || (__      | (__    | (_____ | |      | (___) || (____)|| (__    
//     | |   |  ___  ||  __)     |  __)   (_____  )| |      |  ___  ||  _____)|  __)   
//     | |   | (   ) || (        | (            ) || |      | (   ) || (      | (      
//     | |   | )   ( || (____/\  | (____/\/\____) || (____/\| )   ( || )      | (____/\
//     )_(   |/     \|(_______/  (_______/\_______)(_______/|/     \||/       (_______/
//                                                                                     
//   _______  _______  _______  _______                                                
//  (  ____ )(  ___  )(  ___  )(       )                                               
//  | (    )|| (   ) || (   ) || () () |                                               
//  | (____)|| |   | || |   | || || || |                                               
//  |     __)| |   | || |   | || |(_)| |                                               
//  | (\ (   | |   | || |   | || |   | |                                               
//  | ) \ \__| (___) || (___) || )   ( |                                               
//  |/   \__/(_______)(_______)|/     \|                                               
//                                         
		`,
		"",
		"Welcome to the Escape Room!",
		"",
		"You can [L]ook Around, [I]nspect objects, [U]se Objects, or [Q]uit.",
	}, "\n")
}

type UserConnection struct {
	conn net.Conn
}

func (u *UserConnection) Send(line string) error {
	_, err := u.conn.Write([]byte(line + "\n"))
	return err
}

func (u *UserConnection) Prompt() error {
	_, err := u.conn.Write([]byte(">"))
	return err
}
