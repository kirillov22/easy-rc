package actions

import "fmt"

type Move struct {
	moveX int32
	moveY int32
}

func (p Move) Process(robot Robot) (R any, err error) {
	fmt.Printf("Processing MOVE! I like to move it move it. Moving the x-axis: %d, the y-axis: %d\n", p.moveX, p.moveY)

	currentX, currentY := robot.Location()
	fmt.Println("Current position:", currentX, currentY)

	newX := currentX + int(p.moveX)
	newY := currentY + int(p.moveY)
	robot.Move(newX, newY)
	fmt.Println("Moved to position:", newX, newY)
	return nil, nil
}

func (p Move) Debug() {
	fmt.Println("MOVE!")
}
