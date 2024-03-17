package ws

import (
	"schiffeversenken/data"
	"schiffeversenken/player"
)

type WebsocketPlayer struct {
	Name    string
	Channel chan player.OutMessage
	Field   *data.Field
}

func NewWebsocketPlayer() WebsocketPlayer {
	channel := make(chan player.OutMessage)
	player := WebsocketPlayer{Channel: channel}
	return player
}

func (c *WebsocketPlayer) GetName() string {
	return c.Name
}

func (c *WebsocketPlayer) SetName(name string) {
	c.Name = name
}

func (c *WebsocketPlayer) GetField() *data.Field {
	return c.Field
}

func (c *WebsocketPlayer) GetChan() chan player.OutMessage {
	return c.Channel
}

func (c *WebsocketPlayer) CreateField() {
	field := data.NewField(data.Vector{X: 10, Y: 10})
	field.Settings = data.GetDefaultFieldSettings()
	c.Field = &field
}
