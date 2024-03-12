package game

import (
	"schiffeversenken/player"
)

var games = make(map[*Game]chan player.InMessage, 0)
var players = make(map[*player.Player]*Game)

func CreateGame(player1 *player.Player, player2 *player.Player) {
	channel := make(chan player.InMessage)
	game := NewGame(player1, player2, channel)
	games[&game] = channel
}

func GetGame(player *player.Player) (*Game, bool) {
	game, exists := players[player]
	return game, exists
}

func GetGameChannel(game *Game) chan player.InMessage {
	return games[game]
}
