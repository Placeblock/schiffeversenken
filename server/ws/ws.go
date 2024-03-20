package ws

import (
	"encoding/json"
	"fmt"
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
	go Listen(&p, con)
	defer con.Close()
	defer close(p.GetChan())
	defer match.RemovePlayer(&p)

	match.AddToPool(&p)
	for {
		_, data, err := con.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		var message WebsocketMessage
		err = json.Unmarshal(data, &message)
		if err != nil {
			fmt.Println(err)
			p.Channel <- player.OutMessage{Action: "INVALID_MESSAGE"}
			continue
		}
		switch message.Action {
		case "POOL":
			match.AddToPool(&p)
		case "JOIN":
			var id string
			err = json.Unmarshal(*message.Data, &id)
			if err != nil {
				continue
			}
			match.Join(&p, id)
		default:
			gameMessage := player.GetGameMessage(message.Action, *message.Data)
			if gameMessage == nil {
				continue
			}
			channel := match.GetGameChannel(&p)
			channel <- player.InMessage{Player: &p, Data: gameMessage}
		}
	}
}

func Listen(p *WebsocketPlayer, con *websocket.Conn) {
	for message := range p.Channel {
		con.WriteJSON(message)
	}
}

func StartServer() {
	http.HandleFunc("/", handle)
	http.ListenAndServe("127.0.0.1:4195", nil)
}
