package actions

type Robot interface {
	Move(x, y int)
	Location() (int, int)
	Click(button string)
}
