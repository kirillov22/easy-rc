package actions

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type Move struct {
	moveX int32
	moveY int32
}

func (p Move) Process() (R any, err error) {
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
