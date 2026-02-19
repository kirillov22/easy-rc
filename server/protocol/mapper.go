package protocol

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"
	"fmt"
	"time"
)

func FromProto(message *proto_messages.Message) (actions.Processable, error) {
	switch msg := message.GetMsg().(type) {
	case *proto_messages.Message_Ping:
		return actions.NewPing(), nil
	case *proto_messages.Message_Pong:
		return actions.NewPong(time.UnixMilli(msg.Pong.Timestamp)), nil
	case *proto_messages.Message_Move:
		return actions.NewMove(msg.Move.MoveX, msg.Move.MoveY), nil
	case *proto_messages.Message_Click:
		return actions.NewClick(actions.MouseButton(msg.Click.MouseButton)), nil
	default:
		return nil, fmt.Errorf("unknown message type: %T", msg)
	}
}

func ToProto(processable actions.Processable) (*proto_messages.Message, error) {
	switch p := processable.(type) {
	case actions.Pong:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Pong{
				Pong: &proto_messages.Pong{Timestamp: p.Timestamp().UnixMilli()},
			},
		}, nil
	case actions.Ping:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Ping{
				Ping: &proto_messages.Ping{Timestamp: time.Now().UnixMilli()},
			},
		}, nil
	case actions.Move:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Move{
				Move: &proto_messages.Move{MoveX: p.MoveX(), MoveY: p.MoveY()},
			},
		}, nil
	case actions.Click:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Click{
				Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton(p.MouseButton())},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown processable type: %T", p)
	}
}
