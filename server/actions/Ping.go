package actions

import "time"

type Ping struct{}

func (p Ping) Process(_ Robot) (R any, err error) {
	return Pong{timestamp: time.Now()}, nil
}

func (p Ping) String() string {
	return "Ping{}"
}
