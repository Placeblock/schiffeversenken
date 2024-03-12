package console

import (
	"schiffeversenken/data"
	"schiffeversenken/player"
)

type Player interface {
	GetName() string
	GetChan() chan player.OutMessage
	GetField() *data.Field
}

type ConsolePlayer struct {
	Name    string
	Channel chan player.OutMessage
	Field   *data.Field
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
