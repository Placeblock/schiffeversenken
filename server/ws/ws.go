package ws

import (
	"encoding/json"
	"log"
	"net/http"
	"schiffeversenken/match"
	"schiffeversenken/player"
	"sync"
	"time"

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
	lock := sync.Mutex{}
	keepAlive(con, time.Second*20, &lock)
	p := NewWebsocketPlayer()
	go Listen(&p, con, &lock)
	defer con.Close()
	defer close(p.GetChan())
	defer match.RemovePlayer(&p)

	match.AddToPool(&p)
	for {
		_, data, err := con.ReadMessage()
		if err != nil {
			break
		}
		var message WebsocketMessage
		err = json.Unmarshal(data, &message)
		if err != nil {
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

func keepAlive(c *websocket.Conn, timeout time.Duration, lock *sync.Mutex) {
	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		for {
			lock.Lock()
			err := c.WriteMessage(websocket.PingMessage, []byte("keepalive"))
			lock.Unlock()
			if err != nil {
				return
			}
			time.Sleep(timeout / 2)
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

func Listen(p *WebsocketPlayer, con *websocket.Conn, lock *sync.Mutex) {
	for message := range p.Channel {
		lock.Lock()
		con.WriteJSON(message)
		lock.Unlock()
	}
}

func StartServer() {
	http.HandleFunc("/", handle)
	http.ListenAndServe("127.0.0.1:4195", nil)
}
