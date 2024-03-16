package player

import (
	"schiffeversenken/data"
)

type Player interface {
	GetName() string
	GetChan() chan OutMessage
	GetField() *data.Field
}

type OutMessage struct {
	Action string
	Data   interface{}
}

type InMessage struct {
	Player Player
	Action string
	Data   interface{}
}

type ShootData struct {
	Cell data.Vector `json:"cell"`
}

type ShipData struct {
	Position  data.Vector `json:"position"`
	Length    uint8       `json:"length"`
	Direction data.Vector `json:"direction"`
}
