package websocket

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"
	"easy-rc-server/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func dialTestServer(t *testing.T, robot *testutil.FakeRobot) (*websocket.Conn, *httptest.Server) {
	t.Helper()
	server := httptest.NewServer(NewServer(robot))
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	if err != nil {
		server.Close()
		t.Fatalf("failed to dial WebSocket: %v", err)
	}
	return conn, server
}

func sendProto(t *testing.T, conn *websocket.Conn, msg *proto_messages.Message) {
	t.Helper()
	data, err := proto.Marshal(msg)
	if err != nil {
		t.Fatalf("failed to marshal message: %v", err)
	}
	if err := conn.WriteMessage(websocket.BinaryMessage, data); err != nil {
		t.Fatalf("failed to write message: %v", err)
	}
}

func readProto(t *testing.T, conn *websocket.Conn) *proto_messages.Message {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))
	_, data, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("failed to read message: %v", err)
	}
	msg := &proto_messages.Message{}
	if err := proto.Unmarshal(data, msg); err != nil {
		t.Fatalf("failed to unmarshal response: %v", err)
	}
	return msg
}

func syncWithPing(t *testing.T, conn *websocket.Conn) {
	t.Helper()
	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Ping{Ping: &proto_messages.Ping{}},
	})
	resp := readProto(t, conn)
	if resp.GetPong() == nil {
		t.Fatal("expected pong response during sync")
	}
}

func TestWSPingReturnsPong(t *testing.T) {
	robot := &testutil.FakeRobot{}
	conn, server := dialTestServer(t, robot)
	defer server.Close()
	defer conn.Close()

	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Ping{
			Ping: &proto_messages.Ping{Timestamp: 12345},
		},
	})

	resp := readProto(t, conn)
	if resp.GetPong() == nil {
		t.Fatal("expected pong response")
	}
}

func TestWSMoveUpdatesRobot(t *testing.T) {
	robot := &testutil.FakeRobot{X: 100, Y: 200}
	conn, server := dialTestServer(t, robot)
	defer server.Close()
	defer conn.Close()

	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Move{
			Move: &proto_messages.Move{MoveX: 50, MoveY: -30},
		},
	})

	syncWithPing(t, conn)

	if len(robot.Moves) != 1 {
		t.Fatalf("expected 1 move, got %d", len(robot.Moves))
	}
	if robot.Moves[0].X != 150 || robot.Moves[0].Y != 170 {
		t.Errorf("expected move to (150, 170), got (%d, %d)", robot.Moves[0].X, robot.Moves[0].Y)
	}
}

func TestWSClickUpdatesRobot(t *testing.T) {
	robot := &testutil.FakeRobot{}
	conn, server := dialTestServer(t, robot)
	defer server.Close()
	defer conn.Close()

	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Click{
			Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_RIGHT},
		},
	})

	syncWithPing(t, conn)

	if len(robot.Clicks) != 1 {
		t.Fatalf("expected 1 click, got %d", len(robot.Clicks))
	}
	if robot.Clicks[0] != actions.RightButton {
		t.Errorf("expected right click, got %s", robot.Clicks[0])
	}
}

func TestWSMultipleMessages(t *testing.T) {
	robot := &testutil.FakeRobot{X: 0, Y: 0}
	conn, server := dialTestServer(t, robot)
	defer server.Close()
	defer conn.Close()

	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Move{
			Move: &proto_messages.Move{MoveX: 10, MoveY: 20},
		},
	})
	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Move{
			Move: &proto_messages.Move{MoveX: 5, MoveY: 5},
		},
	})
	sendProto(t, conn, &proto_messages.Message{
		Msg: &proto_messages.Message_Click{
			Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_LEFT},
		},
	})

	syncWithPing(t, conn)

	if len(robot.Moves) != 2 {
		t.Fatalf("expected 2 moves, got %d", len(robot.Moves))
	}
	if len(robot.Clicks) != 1 {
		t.Fatalf("expected 1 click, got %d", len(robot.Clicks))
	}
}

func TestWSInvalidMessageDoesNotKillConnection(t *testing.T) {
	robot := &testutil.FakeRobot{}
	conn, server := dialTestServer(t, robot)
	defer server.Close()
	defer conn.Close()

	// Send garbage bytes
	if err := conn.WriteMessage(websocket.BinaryMessage, []byte{0xFF, 0xFE, 0xFD}); err != nil {
		t.Fatalf("failed to write garbage: %v", err)
	}

	// Connection should still be alive after garbage
	syncWithPing(t, conn)
}

func TestWSUpgradeRejectsNonLocalOrigin(t *testing.T) {
	robot := &testutil.FakeRobot{}
	server := httptest.NewServer(NewServer(robot))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	dialer := websocket.Dialer{}
	header := http.Header{}
	header.Set("Origin", "http://8.8.8.8")

	_, resp, err := dialer.Dial(wsURL, header)
	if err == nil {
		t.Fatal("expected connection to be rejected for non-local origin")
	}
	if resp != nil && resp.StatusCode != http.StatusForbidden {
		t.Logf("connection rejected with status %d (expected 403)", resp.StatusCode)
	}
}
