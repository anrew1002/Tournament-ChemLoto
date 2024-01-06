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
	Name       string         `json:"roomName" validate:"required"`
	Time       int            `json:"time" validate:"required"`
	Max_partic int            `json:"maxPlayers" validate:"required"`
	Elements   map[string]int `json:"elementCounts" validate:"required"`
	IsAuto     bool           `json:"isAuto"`
}
