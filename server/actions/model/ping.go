package model

import (
	"easy-rc-server/actions"
	"time"
)

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

func (p Ping) Process(_ actions.Robot) (actions.Processable, error) {
	return Pong{timestamp: time.Now()}, nil
}

func (p Ping) String() string {
	return "Ping{}"
}
