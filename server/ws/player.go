package ws

import (
	"schiffeversenken/data"
	"schiffeversenken/player"
)

type WebsocketPlayer struct {
	Channel chan player.OutMessage
	Field   *data.Field
}

func NewWebsocketPlayer() WebsocketPlayer {
	channel := make(chan player.OutMessage)
	player := WebsocketPlayer{Channel: channel}
	return player
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
	c.Channel <- player.OutMessage{Action: "FIELD", Data: c.Field}
}
