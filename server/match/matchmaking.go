package match

import (
	"schiffeversenken/game"
	"schiffeversenken/player"
)

var players = make(map[string]player.Player)

var games = make(map[player.Player]*game.Game)
var gameChannels = make(map[*game.Game]chan player.InMessage)

func CreateGame(player1 player.Player, player2 player.Player) {
	delete(players, player1.GetName())
	delete(players, player2.GetName())
	channel := make(chan player.InMessage)
	g := game.NewGame(player1, player2, channel)
	go g.Listen()
	games[player1] = &g
	games[player2] = &g
	gameChannels[&g] = channel
}

func GetGameChannel(p player.Player) chan player.InMessage {
	game := games[p]
	if game == nil {
		return nil
	}
	return game.Channel
}

func getPlayer(name string) player.Player {
	return players[name]
}

func addToPool(p player.Player) {
	players[p.GetName()] = p
}

func CheckName(p player.Player, name string) {
	if nameExists(name) {
		p.GetChan() <- player.OutMessage{Action: "NAME", Data: "ALREADY_EXISTS"}
		return
	}
	p.SetName(name)
	addToPool(p)
}

func CheckOpponent(p player.Player, opponentName string) {
	if opponentName == p.GetName() {
		p.GetChan() <- player.OutMessage{Action: "PLAY", Data: "INVALID_PLAYER"}
		return
	}
	opponent := getPlayer(opponentName)
	if opponent == nil {
		p.GetChan() <- player.OutMessage{Action: "PLAY", Data: "INVALID_PLAYER"}
		return
	}
	CreateGame(p, opponent)
}

func nameExists(name string) bool {
	_, exists := players[name]
	return exists
}
