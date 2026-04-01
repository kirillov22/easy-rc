package model

import "easy-rc-server/actions"

type Processable interface {
	Process(robot actions.Robot) (Processable, error)
}
