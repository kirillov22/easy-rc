package actions

type MouseButton uint16

const (
	LeftButton MouseButton = iota
	MiddleButton
	RightButton
)

func (m MouseButton) String() string {
	switch m {
	case LeftButton:
		return "left"
	case RightButton:
		return "right"
	case MiddleButton:
		return "middle"
	default:
		return "left"
	}
}
