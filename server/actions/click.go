package actions

import "fmt"

type Click struct {
	mouseButton MouseButton
}

func NewClick(button MouseButton) Click {
	return Click{mouseButton: button}
}

func (c Click) MouseButton() MouseButton {
	return c.mouseButton
}

func (c Click) Process(robot Robot) (Processable, error) {
	robot.Click(c.mouseButton)
	return nil, nil
}

func (c Click) String() string {
	return fmt.Sprintf("Click{button=%s}", c.mouseButton)
}
