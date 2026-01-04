package actions

import (
	"fmt"
	"time"
)

type Ping struct {
}

func (p Ping) Process() (R any, err error) {
	pong := Pong{timestamp: time.Now()}
	fmt.Printf("Processing ping! Sending back a pong: %s\n", pong)
	return pong, nil
}

func (p Ping) Debug() {
	fmt.Println("PING!")
}
