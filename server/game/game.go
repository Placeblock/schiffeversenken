package game

import (
	"fmt"
	"schiffeversenken/data"
	"schiffeversenken/player"
)

const (
	BUILDING = iota
	PLAYING  = iota
)

type Game struct {
	State         uint
	Channel       chan player.InMessage
	Player1       *player.Player
	Player2       *player.Player
	CurrentPlayer *player.Player
}

func NewGame(player1 *player.Player, player2 *player.Player, channel chan player.InMessage) Game {
	game := Game{State: BUILDING, Player1: player1, Player2: player2, CurrentPlayer: player1, Channel: channel}
	go game.Listen()
	return game
}

func (g *Game) End() {
	close(g.Channel) // Stop listening for messages
}

func (g *Game) RemovePlayer(removedPlayer *player.Player) {
	g.broadcast("REMOVE_PLAYER", (*removedPlayer).GetName())
	if removedPlayer == g.Player1 {
		g.broadcast("WIN", (*g.Player2).GetName())
	} else {
		g.broadcast("WIN", (*g.Player1).GetName())
	}
}

func (g *Game) getOtherPlayer() *player.Player {
	if g.Player1 == g.CurrentPlayer {
		return g.Player2
	}
	return g.Player1
}

func (g *Game) nextPlayer() {
	g.CurrentPlayer = g.getOtherPlayer()
	g.broadcast("CURRENT_PLAYER", (*g.CurrentPlayer).GetName())
}

func (g *Game) broadcast(action string, data interface{}) {
	(*g.Player1).GetChan() <- player.OutMessage{Action: action, Data: data}
	(*g.Player2).GetChan() <- player.OutMessage{Action: action, Data: data}
}

func (g *Game) PlaceShip(pl *player.Player, ship data.Ship) {
	if !(*pl).GetField().CanAddShip(&ship) {
		(*pl).GetChan() <- player.OutMessage{Action: "INVALID_SHIP", Data: nil}
		return
	}
	(*pl).GetField().AddShip(&ship)
	(*pl).GetChan() <- player.OutMessage{Action: "SHIP_PLACED", Data: nil}
}

type HitResponse struct {
	Player string      `json:"player"`
	Cell   data.Vector `json:"cell"`
}

type SunkResponse struct {
	Player string    `json:"player"`
	Ship   data.Ship `json:"ship"`
}

func (g *Game) Shoot(pl *player.Player, cell data.Vector) {
	if g.State != PLAYING {
		(*pl).GetChan() <- player.OutMessage{Action: "INVALID_STATE", Data: nil}
		return
	}
	if g.CurrentPlayer != pl {
		(*pl).GetChan() <- player.OutMessage{Action: "INVALID_TURN", Data: nil}
		return
	}
	target := g.getOtherPlayer()
	if !(*target).GetField().CanShoot(cell) {
		(*pl).GetChan() <- player.OutMessage{Action: "INVALID_SHOT", Data: nil}
		return
	}
	hit, sunk := (*target).GetField().Shoot(cell)
	if hit {
		g.broadcast("HIT", HitResponse{Player: (*target).GetName(), Cell: cell})

		if sunk {
			ship := (*target).GetField().Cells[cell].Ship
			g.broadcast("SUNK", SunkResponse{Player: (*target).GetName(), Ship: *ship})

			if (*target).GetField().IsDefeated() {
				g.broadcast("DEFEAT", (*target).GetName())
			}
		}
	} else {
		(*pl).GetChan() <- player.OutMessage{Action: "NO_HIT", Data: nil}
	}
	g.nextPlayer()
}

func (g *Game) Listen() {
	for message := range g.Channel {
		switch action := message.Action; action {
		case "SHOOT":
			shootData := message.Data.(player.ShootData)
			g.Shoot(message.Player, shootData.Cell)
		case "ADD_SHIP":
			shipData := message.Data.(player.ShipData)
			ship := data.NewShip(shipData.Position, shipData.Direction, shipData.Length)
			g.PlaceShip(message.Player, ship)
		default:
			fmt.Println("Invalid Message: ", action, message)
		}
	}
}
