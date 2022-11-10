package main

import (
	"math/rand"
	"strconv"
)

// Position describes two things: the room that an entity exists and the location in
// that room. A room can have many locations, and there is a default location in the
// center of the room. This location is the place where items are dropped or players
// exist when they first enter a room. Note that a player can be at a room and a location.
// Players may not coexist in the same location, except the center of the room
type Position struct {
	// Reports the name of the room the object is in. Only set if the object
	// is not in a player's inventory.
	Room     string
	Location string
	// These fields are set for an item that is in someone's inventory.
	InInventory      bool
	PlayersInventory string
}

type Player struct {
	conn *UserConnection
	Name string

	Position
}

func NewPlayer(userConn *UserConnection, name string) *Player {
	p := Player{
		conn: userConn,
		Name: strconv.Itoa(rand.Int()),
		Position: Position{
			Room:     "first room",
			Location: "", // Position is empty since they start in the center of the room
		},
	}
	p.Name = name
	p.conn.playerID = p.Name
	p.conn.start()
	return &p
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

/*
 */

/*
 */
