package actions

import (
	"fmt"
	"time"

	proto_messages "easy-rc-server/generated/proto-messages"
)

func FromProto(message *proto_messages.Message) (Processable, error) {
	switch msg := message.GetMsg().(type) {
	case *proto_messages.Message_Ping:
		return Ping{}, nil
	case *proto_messages.Message_Pong:
		return Pong{timestamp: time.UnixMilli(msg.Pong.Timestamp)}, nil
	case *proto_messages.Message_Move:
		return Move{moveX: msg.Move.MoveX, moveY: msg.Move.MoveY}, nil
	case *proto_messages.Message_Click:
		return Click{mouseButton: MouseButton(msg.Click.MouseButton)}, nil
	default:
		return nil, fmt.Errorf("unknown message type: %T", msg)
	}
}

func ToProto(processable Processable) (*proto_messages.Message, error) {
	switch p := processable.(type) {
	case Pong:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Pong{
				Pong: &proto_messages.Pong{Timestamp: p.timestamp.UnixMilli()},
			},
		}, nil
	case Ping:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Ping{
				Ping: &proto_messages.Ping{Timestamp: time.Now().UnixMilli()},
			},
		}, nil
	case Move:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Move{
				Move: &proto_messages.Move{MoveX: p.moveX, MoveY: p.moveY},
			},
		}, nil
	case Click:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Click{
				Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton(p.mouseButton)},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown processable type: %T", p)
	}
}
