package model

import (
	"easy-rc-server/actions"
	"fmt"
)

type Pong struct{}

func NewPong() Pong {
	return Pong{}
}

func (p Pong) Process(_ actions.Robot) (Processable, error) {
	return nil, nil
}

func (p Pong) String() string {
	return fmt.Sprintf("Pong{}")
}
