package actions

import "fmt"

type Click struct {
	mouseButton MouseButton
}

func (c Click) Process(robot Robot) (R any, err error) {
	fmt.Printf("Processing CLICK! Click, click, clack. %v\n", c.mouseButton)

	robot.Click(c.mouseButton.String())
	return nil, nil
}

func (c Click) Debug() {
	fmt.Println("CLICK!")
}
