package actions

import (
	"fmt"
	"time"
)

type Pong struct {
	timestamp time.Time
}

func (p Pong) Process() (R any, err error) {
	fmt.Println("Processing pong! Doing nothing extra")
	return nil, nil
}

func (p Pong) Debug() {
	fmt.Println("PONG!")
}
