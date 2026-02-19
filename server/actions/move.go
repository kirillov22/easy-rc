package actions

import "fmt"

type Move struct {
	moveX int32
	moveY int32
}

func NewMove(x, y int32) Move {
	return Move{moveX: x, moveY: y}
}

func (m Move) MoveX() int32 {
	return m.moveX
}

func (m Move) MoveY() int32 {
	return m.moveY
}

func (m Move) Process(robot Robot) (Processable, error) {
	currentX, currentY := robot.Location()
	newX := currentX + int(m.moveX)
	newY := currentY + int(m.moveY)
	robot.Move(newX, newY)
	return nil, nil
}

func (m Move) String() string {
	return fmt.Sprintf("Move{x=%d, y=%d}", m.moveX, m.moveY)
}
