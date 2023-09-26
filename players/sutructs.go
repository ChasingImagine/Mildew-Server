package players

import "aftermildewserver/transforms"

type Player struct {
	Id         string                `json:"id"`
	Transforms transforms.Transforms `json:"transforms"`
	PlayerType int                   `json:"playerTypes"`
}
