package match

import (
	"schiffeversenken/game"
	"schiffeversenken/player"
)

var players map[string]player.Player

var games map[player.Player]*game.Game
var gameChannels map[*game.Game]chan player.InMessage

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

func GetPlayer(name string) player.Player {
	return players[name]
}

func AddToPool(p player.Player) {
	players[p.GetName()] = p
}

func NameExists(name string) bool {
	_, exists := players[name]
	return exists
}
