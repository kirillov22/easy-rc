package messages

import (
	"fmt"
)

type Ping struct {
	Packet *Packet
}

func NewPing() Processable {
	return Ping{
		&Packet{
			Header:  1,
			Payload: make([]byte, 8),
		},
	}
}

func (p Ping) Process() (returnPayload *Packet, err error) {
	fmt.Println("Processing ping! Sending back a pong")
	return nil, nil
}

func (p Ping) Debug() {
	fmt.Println("PING!")
}
