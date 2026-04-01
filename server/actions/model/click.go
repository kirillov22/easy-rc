package model

import (
	"easy-rc-server/actions"
	"fmt"
)

type Click struct {
	mouseButton actions.MouseButton
}

func NewClick(button actions.MouseButton) Click {
	return Click{mouseButton: button}
}

func (c Click) MouseButton() actions.MouseButton {
	return c.mouseButton
}

func (c Click) Process(robot actions.Robot) (Processable, error) {
	robot.Click(c.mouseButton)
	return nil, nil
}

func (c Click) String() string {
	return fmt.Sprintf("Click{button=%s}", c.mouseButton)
}
