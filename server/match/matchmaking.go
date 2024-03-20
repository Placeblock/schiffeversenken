package match

import (
	"fmt"
	"math/rand"
	"schiffeversenken/game"
	"schiffeversenken/player"
)

var pool = make(map[string]player.Player)

var games = make(map[player.Player]*game.Game)
var gameChannels = make(map[*game.Game]chan player.InMessage)

func createGame(id string, player1 player.Player, player2 player.Player) {
	delete(pool, id)
	for k, v := range pool {
		if v == player1 || v == player2 {
			delete(pool, k)
		}
	}
	channel := make(chan player.InMessage)
	g := game.NewGame(player1, player2, channel)
	go g.Listen()
	games[player1] = &g
	games[player2] = &g
	gameChannels[&g] = channel
}

func getGame(player player.Player) *game.Game {
	return games[player]
}

func GetGameChannel(p player.Player) chan player.InMessage {
	game := games[p]
	if game == nil {
		return nil
	}
	return game.Channel
}

func generateID() string {
	for {
		id := fmt.Sprintf("%04d", rand.Intn(10000))
		if !exists(id) {
			return id
		}
	}
}

func exists(id string) bool {
	_, exists := pool[id]
	return exists
}

func RemovePlayer(p player.Player) {
	for k, v := range pool {
		if v == p {
			delete(pool, k)
		}
	}
	game := getGame(p)
	if game != nil {
		game.RemovePlayer(p)
		delete(games, game.Player1)
		delete(games, game.Player2)
		delete(gameChannels, game)
	}
}

func AddToPool(p player.Player) {
	RemovePlayer(p)
	id := generateID()
	pool[id] = p
	p.GetChan() <- player.OutMessage{Action: "ROOM", Data: id}
}

func Join(p player.Player, id string) {
	opponent, exists := pool[id]
	if !exists {
		p.GetChan() <- player.OutMessage{Action: "INVALID_ROOM", Data: nil}
		return
	}
	createGame(id, p, opponent)
}
