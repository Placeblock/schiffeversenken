package game

import (
	"schiffeversenken/data"
	"schiffeversenken/player"
)

const (
	BUILDING = iota
	PLAYING  = iota
)

type GamePlayer struct {
	Player player.Player
	Field  data.Field
}

type Game struct {
	Players       []*GamePlayer
	CurrentPlayer uint
}

func (g *Game) GetCurrentPlayer() *GamePlayer {
	return g.Players[g.CurrentPlayer]
}

func (g *Game) NextPlayer() {
	g.CurrentPlayer = (g.CurrentPlayer + 1) % uint(len(g.Players))
}

func (g *Game) Shoot(player player.Player) {

}
