package models

type User struct {
	// Id       string
	Username string
	Score    int
	Room     string
	// Password string
	Admin bool
}
type Room struct {
	Name       string
	Time       int
	Max_partic int
	Elements   map[string]int
}
