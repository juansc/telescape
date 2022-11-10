package main

import (
	"fmt"
	"strings"
)

type Event struct {
	EventType string // What the messages is

	Source   string   // Source is the name of the message source. It can be players, entities, or anonymous
	Audience []string // Audience is the name of the audience. If it's empty it goes to everyone
	Message  string

	Requester string // Person requesting to see room
}

// parseEvent takes in a string and parses the event.
func parseEvent(in string) (Event, error) {
	if len(in) > 50 {
		return Event{}, fmt.Errorf("your message is too long")
	}
	parts := strings.SplitN(in, " ", 2)
	fmt.Println(parts, len(parts))
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
	case "w", "who":
		return Event{
			EventType: "who",
		}, nil
	case "say":
		parts = strings.SplitN(parts[1], " ", 2)
		to := parts[0]
		msg := parts[1]
		fmt.Println(parts)
		if len(parts) < 2 {
			return Event{}, fmt.Errorf("Usage: say [player|all] message")
		}
		audience := []string{}
		if strings.ToLower(to) != "all" {
			audience = append(audience, to)
		}
		return Event{
			EventType: "talk",
			Audience:  audience,
			Message:   msg,
		}, nil
	default:
		return Event{}, fmt.Errorf("I don't know what '%s' is", in)
	}
}
