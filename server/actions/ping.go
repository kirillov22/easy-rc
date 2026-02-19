package actions

import "time"

type Ping struct{}

func NewPing() Ping {
	return Ping{}
}

func (p Ping) Process(_ Robot) (Processable, error) {
	return Pong{timestamp: time.Now()}, nil
}

func (p Ping) String() string {
	return "Ping{}"
}
