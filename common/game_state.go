package common

import "sync"

type Location struct {
	X int
	Y int
}

type Cat struct {
	Username string
	Location Location
}

type GameWorld struct {
	Mutex sync.Mutex
	Cats  []Cat
}
