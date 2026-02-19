package actions

import "fmt"

type Click struct {
	mouseButton MouseButton
}

func (c Click) Process(robot Robot) (R any, err error) {
	robot.Click(c.mouseButton.String())
	return nil, nil
}

func (c Click) String() string {
	return fmt.Sprintf("Click{button=%s}", c.mouseButton)
}
