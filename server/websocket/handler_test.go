package websocket

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"
	"easy-rc-server/internal/testutil"
	"testing"

	"google.golang.org/protobuf/proto"
)

func marshalMessage(t *testing.T, msg *proto_messages.Message) []byte {
	t.Helper()
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}
	return data
}

func TestPingReturnsPong(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Ping{
			Ping: &proto_messages.Ping{Timestamp: 12345},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected non-nil response for ping")
	}

	respMsg := &proto_messages.Message{}
	if err := proto.Unmarshal(response, respMsg); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}

	pong := respMsg.GetPong()
	if pong == nil {
		t.Fatal("expected pong message in response")
	}

	if len(robot.Clicks) != 0 || len(robot.Moves) != 0 {
		t.Error("ping should not trigger any robot calls")
	}
}

func TestPongReturnsNothing(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Pong{
			Pong: &proto_messages.Pong{Timestamp: 12345},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response for pong")
	}

	if len(robot.Clicks) != 0 || len(robot.Moves) != 0 {
		t.Error("pong should not trigger any robot calls")
	}
}

func TestClickLeft(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Click{
			Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_LEFT},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response for click")
	}

	if len(robot.Clicks) != 1 || robot.Clicks[0] != actions.LeftButton {
		t.Errorf("expected single left click, got %v", robot.Clicks)
	}
}

func TestClickRight(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Click{
			Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_RIGHT},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response for click")
	}

	if len(robot.Clicks) != 1 || robot.Clicks[0] != actions.RightButton {
		t.Errorf("expected single right click, got %v", robot.Clicks)
	}
}

func TestClickMiddle(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Click{
			Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_MIDDLE},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response for click")
	}

	if len(robot.Clicks) != 1 || robot.Clicks[0] != actions.MiddleButton {
		t.Errorf("expected single middle click, got %v", robot.Clicks)
	}
}

func TestMoveRelative(t *testing.T) {
	robot := &testutil.FakeRobot{X: 200, Y: 300}
	handler := NewMessageHandler(robot)

	data := marshalMessage(t, &proto_messages.Message{
		Msg: &proto_messages.Message_Move{
			Move: &proto_messages.Move{MoveX: 100, MoveY: 50},
		},
	})

	response, err := handler.HandleMessage(data)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response != nil {
		t.Fatal("expected nil response for move")
	}

	if len(robot.Moves) != 1 {
		t.Fatalf("expected 1 move call, got %d", len(robot.Moves))
	}
	if robot.Moves[0].X != 300 || robot.Moves[0].Y != 350 {
		t.Errorf("expected move to (300, 350), got (%d, %d)", robot.Moves[0].X, robot.Moves[0].Y)
	}
}

func TestInvalidBytes(t *testing.T) {
	robot := &testutil.FakeRobot{}
	handler := NewMessageHandler(robot)

	_, err := handler.HandleMessage([]byte{0xFF, 0xFE, 0xFD, 0xFC})
	if err == nil {
		t.Fatal("expected error for invalid bytes")
	}

	if len(robot.Clicks) != 0 || len(robot.Moves) != 0 {
		t.Error("invalid message should not trigger any robot calls")
	}
}

func TestShouldFail(t *testing.T) {
	if 1 != 2 {
		t.Error("invalid message should not trigger any robot calls")
	}
}
