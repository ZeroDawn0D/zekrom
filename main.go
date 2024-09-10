package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

// To serve a directory on disk (/tmp) under an alternate URL
// path (/tmpfiles/), use StripPrefix to modify the request
// URL's path before the FileServer sees it:
// http.Handle("/tmpfiles/", http.StripPrefix("/tmpfiles/", http.FileServer(http.Dir("/tmp"))))

// http.Handle("/app/", http.FileServer(http.Dir("./dist")))
// http://localhost:8080/app/ would try to find a directory named app inside ./dist/, and
// because index.html is not within an app subdirectory inside dist, the request would result in a 404 error.

/*
/turn-based

/turn-based/http

/turn-based/http/server

/turn-based/http/server/state

/turn-based/http/players

/turn-based/http/players/1
/turn-based/http/players/2
*/

func main() {
	state := GameState{Awaiting_action: "p1"}
	origin := Position{0, 0}
	state.Players = []Position{origin, origin}

	r := mux.NewRouter()
	err := GetNetworkIP()
	if err != nil {
		return
	}
	//serve HandleVue before as defineRoutes has a more generic structure, which triggers the js css files route
	HandleVue(r)
	defineRoutes(r, &state)

	//FixMimeTypes()

	HandleHTTPServer(r)
}

func PlayerAction(w http.ResponseWriter, r *http.Request, state *GameState, action ActionType, player uint8) {
	//state.player
	player_ind := player - 1
	curPosition := state.Players[player_ind]
	CurX := curPosition.X
	CurY := curPosition.Y
	switch action {
	case LEFT:
		if CurX > 0 {
			curPosition = Position{X: CurX - 1, Y: CurY}
		}
	case RIGHT:
		if CurX < 5 {
			curPosition = Position{X: CurX + 1, Y: CurY}
		}
	case UP:
		if CurY > 0 {
			curPosition = Position{X: CurX, Y: CurY - 1}
		}
	case DOWN:
		if CurY < 5 {
			curPosition = Position{X: CurX, Y: CurY + 1}
		}
	}
	state.Players[player_ind] = curPosition
	err := json.NewEncoder(w).Encode(state)
	if err != nil {
		panic(err)
	}
}
