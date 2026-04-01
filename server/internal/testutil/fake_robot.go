package testutil

import "easy-rc-server/actions"

type FakeRobot struct {
	X, Y   int
	Clicks []actions.MouseButton
	Moves  []struct{ X, Y int }
}

func (f *FakeRobot) Move(x, y int) {
	f.X = x
	f.Y = y
	f.Moves = append(f.Moves, struct{ X, Y int }{x, y})
}

func (f *FakeRobot) Location() (int, int) {
	return f.X, f.Y
}

func (f *FakeRobot) Click(button actions.MouseButton) {
	f.Clicks = append(f.Clicks, button)
}
