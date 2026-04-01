package model

import (
	"easy-rc-server/actions"
)

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

func (p Ping) Process(_ actions.Robot) (Processable, error) {
	return Pong{}, nil
}

func (p Ping) String() string {
	return "Ping{}"
}
