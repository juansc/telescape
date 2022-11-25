package main

import (
	"fmt"
	"sync"
	"time"
)

type nameGen struct {
	sync.Mutex
	Names   []string
	counter int
}

func NewNameGen(names []string) nameGen {
	return nameGen{Names: names}
}

func (n *nameGen) getName() string {
	n.Lock()
	name := n.Names[n.counter]
	n.counter++
	n.Unlock()
	return name
}

type Engine struct {
	players    []*Player
	playerByID map[string]*Player
	time.Ticker
	events []Event
	close  chan struct{}
	mu     sync.Mutex
}

func NewEngine() *Engine {
	engine := Engine{
		players: []*Player{},
		Ticker:  *time.NewTicker(time.Millisecond * 10),
		events:  []Event{},
		close:   make(chan struct{}),
	}
	go engine.Tick()
	return &engine
}

func (engine *Engine) AddPlayer(p *Player) {
	engine.mu.Lock()
	engine.players = append(engine.players, p)
	engine.mu.Unlock()
}

func (engine *Engine) Tick() {
	for {
		select {
		case <-engine.close:
			return
		case <-engine.Ticker.C:
			engine.mu.Lock()
			// handle and resolve all events
			for i := range engine.events {
				engine.handleEvent(engine.events[i])
			}
			// flush messages for all players
			for i := range engine.players {
				engine.players[i].conn.FlushMessages()
			}
			// empty out the events for the next loop
			engine.events = []Event{}
			engine.mu.Unlock()
		}
	}
}

func (engine *Engine) addEvent(event Event) {
	engine.mu.Lock()
	engine.events = append(engine.events, event)
	engine.mu.Unlock()
}

var playerColors = map[string]string{
	RoleArchitect: "yellow",
	RoleExplorer:  "green",
	RoleThief:     "blue",
	RoleCharlatan: "magenta",
}

const (
	RoleExplorer  = "Explorer"
	RoleArchitect = "Architect"
	RoleThief     = "Thief"
	RoleCharlatan = "Charlatan"
)

func (engine *Engine) handleEvent(event Event) {
	if event.EventType == "talk" {
		audienceMap := map[string]struct{}{}
		for _, p := range event.Audience {
			audienceMap[p] = struct{}{}
		}
		toAll := len(event.Audience) == 0
		fmt.Println("is this going to everyone? ", toAll)
		fmt.Println("the audience map is ", audienceMap)
		msg := fmt.Sprintf("\n[%s]: %s", event.Source, event.Message)
		for i, p := range engine.players {
			if _, ok := audienceMap[p.conn.playerID]; ok || toAll && p.conn.playerID != event.Source {
				// TODO: Make the color specific to the source. For example, all messages originating from
				// the game are white. All messages coming from the Architect are yellow.
				engine.players[i].conn.SendColor(msg, playerColors[event.Source])
			}
		}

		return
	}

	if event.EventType == "look" {
		for i, p := range engine.players {
			if p.Name == event.Requester {
				engine.players[i].conn.Send(NewFountainRoom().Describe())
				return
			}
		}
	}
}
