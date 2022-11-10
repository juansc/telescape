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

type UserConnection struct {
	conn     net.Conn
	playerID string
	close    chan struct{}
	messages chan string
	isClosed bool
	engine   *Engine
}

func NewUserConnection(conn net.Conn, engine *Engine) (*UserConnection, error) {
	userConn := UserConnection{
		conn:     conn,
		close:    make(chan struct{}, 1),
		messages: make(chan string, 100),
		engine:   engine,
	}
	return &userConn, nil
}

func (u *UserConnection) start() {
	go u.loop()
	go u.msgLoop()
}

// msgLoop writes messages to the user witht their connection.
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
	u.SendColor(IntroString())
	u.Send("Hello, " + u.playerID)
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
		event, err := parseEvent(data)
		if err != nil {
			u.messages <- err.Error()
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

func (u *UserConnection) Send(line string) error {
	_, err := u.conn.Write([]byte(line + "\n"))
	return err
}

func (u *UserConnection) SendColor(line string) error {
	line = fmt.Sprintf("%s%s%s", "\x1b[31m", line, "\x1b[0m")
	_, err := u.conn.Write([]byte(line + "\n"))
	return err
}

func (u *UserConnection) Prompt() error {
	_, err := u.conn.Write([]byte(">"))
	return err
}
