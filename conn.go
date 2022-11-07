package main

import (
	"net"
	"sync"
	"time"
)

type playerConn struct {
	// The connection used to communicate with a player. This
	// is where all the messages are
	tcp *net.Conn
}

type Player struct {
	conn playerConn
	Name string
	Room string
}

// When a person logs in, their connection is treated as a SESSION. That session
// is the way that a player connects to the game.

// TASK: In the game object, add a series of connections so that a TCP connection
// goes in we create a player with a name and a unique ID. Then, when that player
// attempts to connect if they provide a name they can log back in.

// We should also be able to broadcase messages to all players, and also be able to
// send messages from player A to player B. That is, messages have:
// * SOURCE - Can be anonymous
// * AUDIENCE - Can be one or more players
// * MESSAGE - Message

// The engine is the other importatnt component.
// It builds a queue of events and resolves them based on priorities.
// The events for now are messages or requests to look around. The messages
// are always run first.

type Engine struct {
	players []Player
	time.Ticker
	events []Event
	close  chan struct{}
	mu     sync.Mutex
}

func NewEngine() *Engine {
	engine := Engine{
		players: []Player{},
		Ticker:  *time.NewTicker(time.Second),
		events:  []Event{},
		close:   make(chan struct{}),
	}
	// go engine.Tick()
	return &engine
}

// receiveMsg receives a message from a connection and does something with it.
// Any parsing failures are returned immediately, but anything that parses
// into well-formed event will be queued.
func (engine *Engine) receiveMsg(msg string) (string, error) {
	return "", nil
}

type Event struct {
	EventType string // What the messages is

	Source   string   // Source is the name of the source
	Audience []string // Audience is the name of the audience. If it's empty it goes to everyone
	Message  string

	Requester string // Person requesting to see room
}

/*
func (engine *Engine) Tick() {
	for {
		select {
		case <-engine.close:
			return
		case <-engine.Ticker.C:
			engine.mu.Lock()
			for i := range engine.events {
				engine.handleEvent(engine.events[i])
			}
			engine.mu.Unlock()
			fmt.Println("Another tick has passed")
		}
	}
}
*/

/*
func (engine *Engine) handleEvent(event Event) {
	if event.EventType == "message" {
		audienceMap := map[string]struct{}{}
		for _, p := range engine.players {
			audienceMap[p.Name] = struct{}{}
		}
		toAll := len(event.Audience) == 0
		for i, p := range engine.players {
			if _, ok := audienceMap[p.Name]; ok || toAll {
				engine.players[i].conn.tcp.Write([]byte(event.Message + "\n"))
			}
		}
		return
	}

	if event.EventType == "room_request" {
		for i, p := range engine.players {
			if p.Name == event.Requester {
				engine.players[i].conn.tcp.Write([]byte("You are in an empty room\n"))
				return
			}
		}
	}
}
*/
