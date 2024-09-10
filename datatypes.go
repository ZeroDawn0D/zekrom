package main

var A string = "working"

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type GameState struct {
	Players         []Position `json:"players"`
	Last_turn_count uint32     `json:"last_turn_count"`
	Awaiting_action string     `json:"awaiting_action"`
}

type ServerQuery struct {
	Action ActionType `json:"action"`
	Player uint8      `json:"player"`
}

type ActionType uint8

const (
	LEFT ActionType = iota
	RIGHT
	UP
	DOWN
	CONNECT
	DISCONNECT
	GET_STATE
)

type serverType uint8

const (
	turnBasedHTTP serverType = iota
	turnBasedWebSocket
	realTimeWebSocket
)

type MultiplayerServer struct {
	st serverType
}
