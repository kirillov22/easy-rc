package messages

import (
	"encoding/binary"
	"fmt"
)

type Packet struct {
	// Header: Defines the message type ie. PING, click, move, etc.
	Header uint16
	// Payload: Contains the info required to process the message
	Payload []byte
}

func FromByteArray(bytes []byte) Packet {
	length := len(bytes)
	if length != 10 {
		fmt.Println("Received unexpected number of bytes!. Expected 10, but got:", length)
	}

	header := binary.BigEndian.Uint16(bytes[0:2])
	payload := bytes[2:10]

	return Packet{Header: header, Payload: payload}
}
