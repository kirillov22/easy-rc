package messages

import (
	"encoding/binary"
	"fmt"

	"github.com/go-vgo/robotgo"
)

type Move struct {
	moveX int32
	moveY int32
}

func NewMove(payload []byte) Processable {
	return Move{
		moveX: int32(binary.BigEndian.Uint32(payload[0:4])),
		moveY: int32(binary.BigEndian.Uint32(payload[4:8])),
	}
}

func (p Move) Process() (returnPayload *Packet, err error) {
	fmt.Printf("Processing MOVE! I like to move it move it. Moving the x-axis: %d, the y-axis: %d\n", p.moveX, p.moveY)

	currentX, currentY := robotgo.Location()
	fmt.Println("Current position:", currentX, currentY)

	newX := currentX + int(p.moveX)
	newY := currentY + int(p.moveY)
	robotgo.Move(newX, newY)
	currentX, currentY = robotgo.Location()
	fmt.Println("Moved to position:", currentX, currentY)
	return nil, nil
}

func (p Move) Debug() {
	fmt.Println("MOVE!")
}
