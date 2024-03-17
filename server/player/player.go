package player

import (
	"encoding/json"
	"schiffeversenken/data"
)

type Player interface {
	GetName() string
	SetName(string)
	GetChan() chan OutMessage
	GetField() *data.Field
	CreateField()
}

type OutMessage struct {
	Action string      `json:"action"`
	Data   interface{} `json:"data"`
}

type InMessage struct {
	Player Player
	Data   interface{} `json:"data"`
}

type ShootData struct {
	Cell data.Vector `json:"cell"`
}

type ShipData struct {
	Position  data.Vector `json:"position"`
	Length    uint8       `json:"length"`
	Direction data.Vector `json:"direction"`
}

func GetGameMessage(action string, data []byte) interface{} {
	switch action {
	case "SHIP":
		var shipData ShipData
		json.Unmarshal(data, &shipData)
		return shipData
	case "SHOOT":
		var shootData ShootData
		json.Unmarshal(data, &shootData)
		return shootData
	default:
		return nil
	}
}
