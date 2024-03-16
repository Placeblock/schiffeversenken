package console

import (
	"fmt"
	"schiffeversenken/data"
	"schiffeversenken/player"
)

type Player interface {
	GetName() string
	GetChan() chan player.OutMessage
	GetGameChan() chan player.InMessage
	GetField() *data.Field
}

type ConsolePlayer struct {
	Name        string
	Channel     chan player.OutMessage
	GameChannel chan player.InMessage
	Field       *data.Field
}

func NewConsolePlayer(name string, gameChannel chan player.InMessage) ConsolePlayer {
	channel := make(chan player.OutMessage)
	field := data.NewField(data.Vector{X: 10, Y: 10})
	player := ConsolePlayer{Name: name, Channel: channel, Field: &field, GameChannel: gameChannel}
	go player.Listen()
	return player
}

func (c *ConsolePlayer) Listen() {
	for message := range c.Channel {
		fmt.Printf("%s: %s -> %+v\n", c.Name, message.Action, message.Data)
	}
}

func (c *ConsolePlayer) GetName() string {
	return c.Name
}

func (c *ConsolePlayer) GetField() *data.Field {
	return c.Field
}

func (c *ConsolePlayer) GetChan() chan player.OutMessage {
	return c.Channel
}

func (c *ConsolePlayer) GetGameChan() chan player.InMessage {
	return c.GameChannel
}
