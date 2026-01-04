package messages

type MouseButton uint16

const (
	LeftButton MouseButton = iota
	RightButton
	MiddleButton
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
