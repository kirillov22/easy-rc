package model

import (
	"easy-rc-server/actions"
	"easy-rc-server/internal/testutil"
	"testing"
)

func TestMoveProcess(t *testing.T) {
	tests := []struct {
		name           string
		startX, startY int
		moveX, moveY   int32
		wantX, wantY   int
	}{
		{"positive deltas", 100, 200, 50, 75, 150, 275},
		{"negative deltas", 100, 200, -30, -50, 70, 150},
		{"zero deltas", 100, 200, 0, 0, 100, 200},
		{"from origin", 0, 0, 10, 20, 10, 20},
		{"negative result", 50, 50, -100, -100, -50, -50},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			robot := &testutil.FakeRobot{X: tt.startX, Y: tt.startY}
			move := NewMove(tt.moveX, tt.moveY)

			result, err := move.Process(robot)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != nil {
				t.Fatal("expected nil result from Move.Process")
			}
			if len(robot.Moves) != 1 {
				t.Fatalf("expected 1 move call, got %d", len(robot.Moves))
			}
			if robot.Moves[0].X != tt.wantX || robot.Moves[0].Y != tt.wantY {
				t.Errorf("expected move to (%d, %d), got (%d, %d)", tt.wantX, tt.wantY, robot.Moves[0].X, robot.Moves[0].Y)
			}
		})
	}
}

func TestClickProcess(t *testing.T) {
	tests := []struct {
		name   string
		button actions.MouseButton
	}{
		{"left", actions.LeftButton},
		{"middle", actions.MiddleButton},
		{"right", actions.RightButton},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			robot := &testutil.FakeRobot{}
			click := NewClick(tt.button)

			result, err := click.Process(robot)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if result != nil {
				t.Fatal("expected nil result from Click.Process")
			}
			if len(robot.Clicks) != 1 {
				t.Fatalf("expected 1 click call, got %d", len(robot.Clicks))
			}
			if robot.Clicks[0] != tt.button {
				t.Errorf("expected %s click, got %s", tt.button, robot.Clicks[0])
			}
		})
	}
}

func TestPingProcess(t *testing.T) {
	robot := &testutil.FakeRobot{}
	ping := NewPing()

	result, err := ping.Process(robot)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, ok := result.(Pong)
	if !ok {
		t.Fatalf("expected Pong, got %T", result)
	}
	if len(robot.Clicks) != 0 || len(robot.Moves) != 0 {
		t.Error("ping should not trigger any robot calls")
	}
}

func TestPongProcess(t *testing.T) {
	robot := &testutil.FakeRobot{}
	pong := NewPong()

	result, err := pong.Process(robot)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != nil {
		t.Fatal("expected nil result from Pong.Process")
	}
	if len(robot.Clicks) != 0 || len(robot.Moves) != 0 {
		t.Error("pong should not trigger any robot calls")
	}
}
