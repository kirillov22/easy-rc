package actions

import (
	"fmt"

	"github.com/go-vgo/robotgo"
)

type Click struct {
	mouseButton MouseButton
}

func (c Click) Process() (R any, err error) {
	fmt.Printf("Processing CLICK! Click, click, clack. %v\n", c.mouseButton)

	robotgo.Click(c.mouseButton.String())
	return nil, nil
}

func (c Click) Debug() {
	fmt.Println("CLICK!")
}
