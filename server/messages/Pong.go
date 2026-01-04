package messages

import (
	"fmt"
)

type Pong struct {
	Packet *Packet
}

func NewPong() Processable {
	return Pong{
		&Packet{
			Header:  2,
			Payload: make([]byte, 8),
		},
	}
}

func (p Pong) Process() (returnPayload *Packet, err error) {
	fmt.Println("Processing pong! Doing nothing extra")
	return p.Packet, nil
}

func (p Pong) Debug() {
	fmt.Println("PONG!")
}
