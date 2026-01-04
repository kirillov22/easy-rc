package messages

import (
	"errors"
	"fmt"
)

func MapMessage(packet *Packet) (Processable, error) {

	if packet == nil {
		return nil, errors.New("packet pointer was nil")
	}

	var message Processable
	switch packet.Header {
	case 1:
		message = NewPing()
		break
	case 2:
		message = NewPong()
		break
	case 3:
		message = NewMove(packet.Payload)
		break
	case 4:
		message = NewClick(packet.Payload)
		break
	default:
		fmt.Printf("WARN: Got unknown packet. Header: %v, payload: %v\n", packet.Header, packet.Payload)
		return nil, errors.New("unexpected packet received")
	}

	return message, nil
}
