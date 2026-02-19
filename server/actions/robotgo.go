package actions

import "github.com/go-vgo/robotgo"

type RobotGo struct{}

func NewRobotGo() *RobotGo {
	return &RobotGo{}
}

func (r *RobotGo) Move(x, y int) {
	robotgo.Move(x, y)
}

func (r *RobotGo) Location() (int, int) {
	return robotgo.Location()
}

func (r *RobotGo) Click(button string) {
	robotgo.Click(button)
}
