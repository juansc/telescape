package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	// Start engine here and start the loop

	// Start the tcp lisener
	listener, err := net.Listen("tcp", "localhost:2000")
	if err != nil {
		fmt.Println("error listening on TCP port", err)
	}

	defer func() {
		listener.Close()
	}()

	fmt.Println("waiting to accept connections...")
	for {
		fmt.Println("waiting for next connection")
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("error listening on accepting connection", err)
		}
		fmt.Println("got a connection")
		_, _ = NewUserConnection(conn)
	}
}

func handleInput(in string) action {
	switch strings.ToLower(in) {
	case "l", "look":
		return action{
			message: "You take a look around your current room. It is empty. You should quit.",
		}
	case "i", "inspect":
		return action{message: "You must specify what you want to inspect."}
	case "u", "use":
		return action{message: "You must specify what you want to use."}
	case "q", "quit":
		return action{message: "Quitter!", do: func() error { return fmt.Errorf("we gotta a quitter") }}
	default:
		return action{message: fmt.Sprintf("I don't know what '%s' is.", in)}
	}
}

type action struct {
	do      func() error
	message string
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
	conn     net.Conn
	close    chan struct{}
	messages chan string
	isClosed bool
}

func NewUserConnection(conn net.Conn) (*UserConnection, error) {
	userConn := UserConnection{
		conn:     conn,
		close:    make(chan struct{}, 1),
		messages: make(chan string, 100),
	}
	go userConn.loop()
	go userConn.msgLoop()
	return &userConn, nil
}

func (u *UserConnection) msgLoop() {
	for {
		select {
		case <-u.close:
			return
		case msg := <-u.messages:
			u.Send(msg)
			u.Prompt()
		}
	}
}

func (u *UserConnection) shutdown() {
	u.isClosed = true
	u.Send("Goodbye!")
	close(u.close)
	u.conn.Close()
}

func (u *UserConnection) loop() {
	u.Send(IntroString())
	u.Prompt()
	defer u.shutdown()
	r := bufio.NewReader(u.conn)
	for {
		u.conn.SetDeadline(time.Now().Add(time.Second * 1))
		data, err := r.ReadString('\n')
		if err != nil {
			if errors.Is(err, os.ErrDeadlineExceeded) {
				continue
			}
			return
		}
		if u.isClosed {
			return
		}
		data = strings.TrimSpace(data)
		// If the user sent only white space then just render the prompt again
		if data == "" {
			u.Prompt()
			continue
		}
		msg, err := parseEvent(data)
		if err != nil {
			u.messages <- err.Error()
			continue
		}
		if msg.EventType == "quit" {
			return
		}
		u.messages <- msg.Message
	}
}

func (u *UserConnection) Send(line string) error {
	_, err := u.conn.Write([]byte(line + "\n"))
	return err
}

func (u *UserConnection) Prompt() error {
	_, err := u.conn.Write([]byte(">"))
	return err
}

// parseEvent takes in a string and parses the event.
func parseEvent(in string) (Event, error) {
	if len(in) > 50 {
		return Event{}, fmt.Errorf("your message is too long")
	}
	parts := strings.Split(in, " ")
	cmd := parts[0]
	switch cmd {
	case "l", "look":
		return Event{
			EventType: "look",
		}, nil
	case "q", "quit":
		return Event{
			EventType: "quit",
		}, nil
	default:
		return Event{}, fmt.Errorf("I don't know what '%s' is", in)
	}
}
