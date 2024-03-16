package console

import (
	"bufio"
	"fmt"
	"os"
	"schiffeversenken/data"
	"schiffeversenken/game"
	"schiffeversenken/player"
	"strconv"
	"strings"
)

func CreateConsoleGame() {
	channel := make(chan player.InMessage)
	player1 := NewConsolePlayer("Felix", channel)
	player2 := NewConsolePlayer("Paula", channel)
	g := game.NewGame(&player1, &player2, channel)
	fmt.Printf("%p\n", &player1)
	fmt.Printf("%p\n", g.Player1)
	go g.Listen()
	go listenConsoleInput(&g)
}

func listenConsoleInput(g *game.Game) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		// convert CRLF to LF
		text = strings.Replace(text, "\n", "", -1)
		contents := strings.Split(text, " ")

		playerName := contents[0]
		var p player.Player
		if g.Player1.GetName() == playerName {
			p = g.Player1
		} else {
			p = g.Player2
		}

		switch contents[1] {
		case "SHIP":
			x, _ := strconv.Atoi(contents[2])
			y, _ := strconv.Atoi(contents[3])
			dx, _ := strconv.Atoi(contents[4])
			dy, _ := strconv.Atoi(contents[5])
			length, _ := strconv.Atoi(contents[6])
			shipData := player.ShipData{Position: data.Vector{X: x, Y: y}, Direction: data.Vector{X: dx, Y: dy}, Length: uint8(length)}
			g.Channel <- player.InMessage{Player: p, Action: "ADD_SHIP", Data: shipData}
		case "SHOOT":
			x, _ := strconv.Atoi(contents[2])
			y, _ := strconv.Atoi(contents[3])
			shootData := player.ShootData{Cell: data.Vector{X: x, Y: y}}
			g.Channel <- player.InMessage{Player: p, Action: "SHOOT", Data: shootData}
		case "START":
			g.Channel <- player.InMessage{Player: p, Action: "START"}
		case "FIELD":
			fmt.Printf("%p\n", p)
			p.GetField().Print()
		}
	}
}
