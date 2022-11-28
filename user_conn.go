package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
	"time"
)

type UserConnection struct {
	conn        net.Conn
	playerID    string
	close       chan struct{}
	isClosed    bool
	engine      *Engine
	mu          sync.Mutex
	outMessages []string
}

func NewUserConnection(conn net.Conn, engine *Engine) (*UserConnection, error) {
	userConn := UserConnection{
		conn:   conn,
		close:  make(chan struct{}, 1),
		engine: engine,
	}
	return &userConn, nil
}

func (u *UserConnection) start() {
	go u.loop()
}

// fastSend immediately sends a message to the user, skipping the buffer that gets flushed periodically.
// It should be for fast fail scenarios like parsing errors.
func (u *UserConnection) fastSend(line string) error {
	u.mu.Lock()
	_, _ = u.conn.Write([]byte(line + "\n"))
	_, _ = u.conn.Write([]byte(">"))
	u.mu.Unlock()
	return nil
}

func (u *UserConnection) FlushMessages() error {
	u.mu.Lock()
	defer u.mu.Unlock()
	// Nothing to do, don't even bother redrawing the prompt.
	if len(u.outMessages) == 0 {
		return nil
	}
	// Send a carriage return to push everything before the person's current buffer
	for _, msg := range u.outMessages {
		// TODO: Handle errors here
		u.send(msg)
	}
	u.outMessages = []string{}
	return nil
}

func (u *UserConnection) shutdown() {
	u.isClosed = true
	u.Send("Goodbye!")
	close(u.close)
	u.conn.Close()
}

func (u *UserConnection) loop() {
	u.SendColor(IntroString(), "red")
	u.Send("Hello, " + u.playerID)
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
			continue
		}
		event, err := parseEvent(data)
		if err != nil {
			u.fastSend(err.Error())
			continue
		}
		event.Source = u.playerID
		event.Requester = u.playerID
		// If the player emitted the quit message then exit the function.
		if event.EventType == "quit" {
			return
		}
		u.engine.addEvent(event)
	}
}

// send a line to the user's TCP connection.
func (u *UserConnection) send(line string) error {
	_, err := u.conn.Write([]byte(line))
	return err
}

// Send
func (u *UserConnection) Send(line string) error {
	u.mu.Lock()
	u.outMessages = append(u.outMessages, line+"\n")
	u.mu.Unlock()
	return nil
}

// SendColor will send the given line with the specified color. If the color does not exist
// it is the same as calling Send(line). The following colors are supported:
// black, red, green, yellow, blue, magenta, cyan, and white.
func (u *UserConnection) SendColor(line string, color string) error {
	// TODO: Add strongly typed strings
	colors := map[string]string{
		"black":   "\x1b[30m",
		"red":     "\x1b[31m",
		"green":   "\x1b[32m",
		"yellow":  "\x1b[33m",
		"blue":    "\x1b[34m",
		"magenta": "\x1b[35m",
		"cyan":    "\x1b[36m",
		"white":   "\x1b[37m",
	}
	line = fmt.Sprintf("%s%s%s", colors[strings.ToLower(color)], line, "\x1b[0m")
	return u.Send(line)
}

// Prompt renders the prompt that the user will see.
func (u *UserConnection) Prompt() error {
	return u.send(">")
}
