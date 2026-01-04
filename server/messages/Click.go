package messages

import (
	"encoding/binary"
	"fmt"

	"github.com/go-vgo/robotgo"
)

type Click struct {
	mouseButton MouseButton
}

func NewClick(payload []byte) Processable {
	intVal := binary.BigEndian.Uint32(payload)
	mouseButton := MouseButton(intVal)

	return Click{
		mouseButton: mouseButton,
	}
}

func (c Click) Process() (returnPayload *Packet, err error) {
	fmt.Printf("Processing CLICK! Click, click, clack. %v\n", c.mouseButton)

	robotgo.Click(c.mouseButton.String())
	return nil, nil
}

func (c Click) Debug() {
	fmt.Println("CLICK!")
}
