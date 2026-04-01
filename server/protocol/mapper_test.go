package protocol

import (
	"easy-rc-server/actions"
	"easy-rc-server/actions/model"
	proto_messages "easy-rc-server/generated/proto-messages"
	"strings"
	"testing"
	"time"
)

func TestFromProto(t *testing.T) {
	t.Run("ping", func(t *testing.T) {
		msg := &proto_messages.Message{
			Msg: &proto_messages.Message_Ping{
				Ping: &proto_messages.Ping{Timestamp: 123},
			},
		}
		result, err := FromProto(msg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if _, ok := result.(model.Ping); !ok {
			t.Fatalf("expected model.Ping, got %T", result)
		}
	})

	t.Run("pong", func(t *testing.T) {
		msg := &proto_messages.Message{
			Msg: &proto_messages.Message_Pong{
				Pong: &proto_messages.Pong{Timestamp: 123},
			},
		}
		result, err := FromProto(msg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		_, ok := result.(model.Pong)
		if !ok {
			t.Fatalf("expected model.Pong, got %T", result)
		}
	})

	t.Run("move", func(t *testing.T) {
		msg := &proto_messages.Message{
			Msg: &proto_messages.Message_Move{
				Move: &proto_messages.Move{MoveX: 10, MoveY: 20},
			},
		}
		result, err := FromProto(msg)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		move, ok := result.(model.Move)
		if !ok {
			t.Fatalf("expected model.Move, got %T", result)
		}
		if move.MoveX() != 10 || move.MoveY() != 20 {
			t.Errorf("expected Move(10, 20), got Move(%d, %d)", move.MoveX(), move.MoveY())
		}
	})

	clickTests := []struct {
		name       string
		protoBtn   proto_messages.MouseButton
		wantButton actions.MouseButton
	}{
		{"click left", proto_messages.MouseButton_MOUSE_BUTTON_LEFT, actions.LeftButton},
		{"click right", proto_messages.MouseButton_MOUSE_BUTTON_RIGHT, actions.RightButton},
		{"click middle", proto_messages.MouseButton_MOUSE_BUTTON_MIDDLE, actions.MiddleButton},
	}
	for _, tt := range clickTests {
		t.Run(tt.name, func(t *testing.T) {
			msg := &proto_messages.Message{
				Msg: &proto_messages.Message_Click{
					Click: &proto_messages.Click{MouseButton: tt.protoBtn},
				},
			}
			result, err := FromProto(msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			click, ok := result.(model.Click)
			if !ok {
				t.Fatalf("expected model.Click, got %T", result)
			}
			if click.MouseButton() != tt.wantButton {
				t.Errorf("expected button %s, got %s", tt.wantButton, click.MouseButton())
			}
		})
	}

	t.Run("nil msg returns error", func(t *testing.T) {
		msg := &proto_messages.Message{}
		_, err := FromProto(msg)
		if err == nil {
			t.Fatal("expected error for nil msg")
		}
		if !strings.Contains(err.Error(), "unknown message type") {
			t.Errorf("expected 'unknown message type' error, got: %v", err)
		}
	})
}

func TestToProto(t *testing.T) {
	t.Run("pong", func(t *testing.T) {
		pong := model.NewPong()
		msg, err := ToProto(pong)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		protoPong := msg.GetPong()
		if protoPong == nil {
			t.Fatal("expected pong message")
		}
	})

	t.Run("ping", func(t *testing.T) {
		ping := model.NewPing()
		before := time.Now().UnixMilli()
		msg, err := ToProto(ping)
		after := time.Now().UnixMilli()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		protoPing := msg.GetPing()
		if protoPing == nil {
			t.Fatal("expected ping message")
		}
		if protoPing.Timestamp < before || protoPing.Timestamp > after {
			t.Errorf("expected timestamp between %d and %d, got %d", before, after, protoPing.Timestamp)
		}
	})

	t.Run("move", func(t *testing.T) {
		move := model.NewMove(5, 10)
		msg, err := ToProto(move)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		protoMove := msg.GetMove()
		if protoMove == nil {
			t.Fatal("expected move message")
		}
		if protoMove.MoveX != 5 || protoMove.MoveY != 10 {
			t.Errorf("expected Move(5, 10), got Move(%d, %d)", protoMove.MoveX, protoMove.MoveY)
		}
	})

	t.Run("click", func(t *testing.T) {
		click := model.NewClick(actions.LeftButton)
		msg, err := ToProto(click)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		protoClick := msg.GetClick()
		if protoClick == nil {
			t.Fatal("expected click message")
		}
		if protoClick.MouseButton != proto_messages.MouseButton_MOUSE_BUTTON_LEFT {
			t.Errorf("expected LEFT, got %v", protoClick.MouseButton)
		}
	})

	t.Run("unknown type returns error", func(t *testing.T) {
		_, err := ToProto(unknownProcessable{})
		if err == nil {
			t.Fatal("expected error for unknown type")
		}
		if !strings.Contains(err.Error(), "unknown processable type") {
			t.Errorf("expected 'unknown processable type' error, got: %v", err)
		}
	})
}

type unknownProcessable struct{}

func (u unknownProcessable) Process(_ actions.Robot) (model.Processable, error) {
	return nil, nil
}
