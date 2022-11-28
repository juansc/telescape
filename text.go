package main

import (
	"fmt"
	"strings"
)

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

type Description string

// A Room has description
type Room struct {
	shortName string
	Description
	Locations []Location
	Items     []Item
}

func (r Room) Describe() string {
	itemStr := ""
	if len(r.Items) > 0 {
		itemStr = "You see the following items:\n"
	}

	for _, i := range r.Items {
		itemStr += fmt.Sprintf("* %s\n", i.shortName)
	}
	return string(describe(
		string(r.Description),
		itemStr,
	))
}

type Location struct {
	Description
	Items []string
}

type Item struct {
	shortName string
	Description
}

func NewLever() Item {
	return Item{
		shortName:   "lever",
		Description: describe("A lever made of milky, smooth marble."),
	}
}

func NewFountainRoom() Room {
	return Room{
		shortName: "fountain_room",
		Description: describe(
			"The room is empty save for a fountain in the middle of the room.",
			"To the north is the only exit.",
		),
		Locations: []Location{NewFountainLocation()},
		Items:     []Item{NewLever()},
	}
}

func NewFountainLocation() Location {
	return Location{
		Description: describe(
			"The fountain is a highly ornate statue of a woman in flowing robes.",
			"Upon closer inspection you can see that the statue is fully articulated",
			"and is capable of movement. You daren't touch it.",
		),
	}
}

func describe(lines ...string) Description {
	return Description(strings.Join(lines, "\n"))
}
