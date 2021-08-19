package common

type Location struct {
	X int
	Y int
}

type Cat struct {
	Username string
	Location Location
}

type GameWorld struct {
	Cats []Cat
}
