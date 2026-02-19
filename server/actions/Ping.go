package actions

import (
	"fmt"
	"time"
)

type Ping struct {
}

func (p Ping) Process(_ Robot) (R any, err error) {
	pong := Pong{timestamp: time.Now()}
	fmt.Printf("Processing ping! Sending back a pong: %v\n", pong)
	return pong, nil
}

func (p Ping) Debug() {
	fmt.Println("PING!")
}
