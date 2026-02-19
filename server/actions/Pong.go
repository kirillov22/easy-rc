package actions

import (
	"fmt"
	"time"
)

type Pong struct {
	timestamp time.Time
}

func (p Pong) Process(_ Robot) (R any, err error) {
	return nil, nil
}

func (p Pong) String() string {
	return fmt.Sprintf("Pong{timestamp=%s}", p.timestamp.Format(time.RFC3339))
}
