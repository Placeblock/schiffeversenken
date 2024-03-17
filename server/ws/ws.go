package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"schiffeversenken/match"
	"schiffeversenken/player"

	"github.com/gorilla/websocket"
)

type WebsocketMessage struct {
	Action string           `json:"action"`
	Data   *json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool {
	return true
}}

func handle(w http.ResponseWriter, r *http.Request) {
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	p := NewWebsocketPlayer()
	defer con.Close()
	for {
		_, data, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var message WebsocketMessage
		err = json.Unmarshal(data, &message)
		if err != nil {
			continue
		}
		if message.Action == "NAME" {
			var name string
			err = json.Unmarshal(*message.Data, &name)

			if err != nil {
				continue
			}

			if match.NameExists(p.Name) {
				p.Channel <- player.OutMessage{Action: "NAME", Data: "ALREADY_EXISTS"}
				continue
			}
			p.Name = name
			match.AddToPool(&p)
			continue
		}
		if message.Action == "PLAY" {
			var opponentName string
			err = json.Unmarshal(*message.Data, &opponentName)

			if err != nil {
				continue
			}

			opponent := match.GetPlayer(opponentName)
			if opponent == nil {
				p.Channel <- player.OutMessage{Action: "PLAY", Data: "INVALID_PLAYER"}
				continue
			}
			match.CreateGame(&p, opponent)
			continue
		}
		channel := match.GetGameChannel(&p)
		if channel == nil {
			continue
		}
		channel <- message
	}
}

func StartServer() {
	http.HandleFunc("/", handle)
	http.ListenAndServe("localhost:4195", nil)
}
