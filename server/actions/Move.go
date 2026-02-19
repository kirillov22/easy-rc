package actions

import "fmt"

type Move struct {
	moveX int32
	moveY int32
}

func (p Move) Process(robot Robot) (R any, err error) {
	currentX, currentY := robot.Location()
	newX := currentX + int(p.moveX)
	newY := currentY + int(p.moveY)
	robot.Move(newX, newY)
	return nil, nil
}

func (p Move) String() string {
	return fmt.Sprintf("Move{x=%d, y=%d}", p.moveX, p.moveY)
}
