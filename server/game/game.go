package game

import (
	"schiffeversenken/data"
	"schiffeversenken/player"
)

const (
	BUILDING = iota
	PLAYING  = iota
	ENDED    = iota
)

type Game struct {
	State         uint
	Channel       chan player.InMessage
	Player1       player.Player
	Player2       player.Player
	CurrentPlayer player.Player
}

func NewGame(player1 player.Player, player2 player.Player, channel chan player.InMessage) Game {
	player1.CreateField()
	player2.CreateField()
	game := Game{State: BUILDING, Player1: player1, Player2: player2, CurrentPlayer: player1, Channel: channel}
	game.broadcast("STATE", "building")
	game.nextPlayer()
	return game
}

func (g *Game) End() {
	g.State = ENDED
	g.broadcast("STATE", "ended")
	close(g.Channel) // Stop listening for messages
}

func (g *Game) RemovePlayer(removedPlayer player.Player) {
	opponent := g.getOtherPlayer(removedPlayer)
	opponent.GetChan() <- player.OutMessage{Action: "OPPONENT_LEFT"}
	g.win(opponent)
}

func (g *Game) getOtherPlayer(p player.Player) player.Player {
	if g.Player1 == p {
		return g.Player2
	}
	return g.Player1
}

func (g *Game) win(p player.Player) {
	p.GetChan() <- player.OutMessage{Action: "WON"}
	g.getOtherPlayer(p).GetChan() <- player.OutMessage{Action: "LOST"}
	g.End()
}

func (g *Game) nextPlayer() {
	g.CurrentPlayer = g.getOtherPlayer(g.CurrentPlayer)
	g.CurrentPlayer.GetChan() <- player.OutMessage{Action: "TURN_START"}
	g.getOtherPlayer(g.CurrentPlayer).GetChan() <- player.OutMessage{Action: "TURN_END"}
}

func (g *Game) broadcast(action string, data interface{}) {
	g.Player1.GetChan() <- player.OutMessage{Action: action, Data: data}
	g.Player2.GetChan() <- player.OutMessage{Action: action, Data: data}
}

func (g *Game) PlaceShip(pl player.Player, ship data.Ship) {
	if g.State != BUILDING {
		pl.GetChan() <- player.OutMessage{Action: "INVALID_STATE", Data: nil}
		return
	}
	if !pl.GetField().CanAddShip(&ship) {
		pl.GetChan() <- player.OutMessage{Action: "INVALID_SHIP", Data: nil}
		return
	}
	pl.GetField().AddShip(&ship)
	pl.GetChan() <- player.OutMessage{Action: "SHIP_PLACED", Data: ship}
	if g.CurrentPlayer.GetField().FinishedPlacing() && g.getOtherPlayer(g.CurrentPlayer).GetField().FinishedPlacing() {
		g.Start()
	}
}

func (g *Game) Start() {
	g.State = PLAYING
	g.broadcast("STATE", "playing")
}

type HitResponse struct {
	Player string      `json:"player"`
	Cell   data.Vector `json:"cell"`
}

type SunkResponse struct {
	Player string    `json:"player"`
	Ship   data.Ship `json:"ship"`
}

func (g *Game) Shoot(pl player.Player, cell data.Vector) {
	if g.State != PLAYING {
		pl.GetChan() <- player.OutMessage{Action: "INVALID_STATE", Data: nil}
		return
	}
	if g.CurrentPlayer != pl {
		pl.GetChan() <- player.OutMessage{Action: "INVALID_TURN", Data: nil}
		return
	}
	target := g.getOtherPlayer(pl)
	if !target.GetField().CanShoot(cell) {
		pl.GetChan() <- player.OutMessage{Action: "INVALID_SHOT", Data: nil}
		return
	}
	hit, sunk := target.GetField().Shoot(cell)
	if hit {
		pl.GetChan() <- player.OutMessage{Action: "HIT_OTHER", Data: cell}
		target.GetChan() <- player.OutMessage{Action: "HIT_SELF", Data: cell}

		if sunk {
			ship := target.GetField().Cells[cell].Ship
			pl.GetChan() <- player.OutMessage{Action: "SUNK_OTHER", Data: *ship}
			target.GetChan() <- player.OutMessage{Action: "SUNK_SELF", Data: *ship}

			if target.GetField().IsDefeated() {
				g.win(pl)
			}
		}
	} else {
		pl.GetChan() <- player.OutMessage{Action: "NO_HIT_OTHER", Data: cell}
		target.GetChan() <- player.OutMessage{Action: "NO_HIT_SELF", Data: cell}
	}
	g.nextPlayer()
}

func (g *Game) Listen() {
	for message := range g.Channel {
		switch payload := message.Data.(type) {
		case player.ShootData:
			g.Shoot(message.Player, payload.Cell)
		case player.ShipData:
			ship := data.NewShip(payload.Position, payload.Direction, payload.Length)
			g.PlaceShip(message.Player, ship)
		}
	}
}
