package actions

import (
	"fmt"
	"time"
)

type Pong struct {
	timestamp time.Time
}

func NewPong(timestamp time.Time) Pong {
	return Pong{timestamp: timestamp}
}

func (p Pong) Timestamp() time.Time {
	return p.timestamp
}

func (p Pong) Process(_ Robot) (Processable, error) {
	return nil, nil
}

func (p Pong) String() string {
	return fmt.Sprintf("Pong{timestamp=%s}", p.timestamp.Format(time.RFC3339))
}
