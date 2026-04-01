package protocol

import (
	"easy-rc-server/actions"
	"easy-rc-server/actions/model"
	proto_messages "easy-rc-server/generated/proto-messages"
	"fmt"
	"time"
)

func FromProto(message *proto_messages.Message) (model.Processable, error) {
	switch msg := message.GetMsg().(type) {
	case *proto_messages.Message_Ping:
		return model.NewPing(), nil
	case *proto_messages.Message_Pong:
		return model.NewPong(), nil
	case *proto_messages.Message_Move:
		return model.NewMove(msg.Move.MoveX, msg.Move.MoveY), nil
	case *proto_messages.Message_Click:
		return model.NewClick(actions.MouseButton(msg.Click.MouseButton)), nil
	default:
		return nil, fmt.Errorf("unknown message type: %T", msg)
	}
}

func ToProto(processable model.Processable) (*proto_messages.Message, error) {
	switch p := processable.(type) {
	case model.Pong:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Pong{
				Pong: &proto_messages.Pong{Timestamp: time.Now().UnixMilli()},
			},
		}, nil
	case model.Ping:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Ping{
				Ping: &proto_messages.Ping{Timestamp: time.Now().UnixMilli()},
			},
		}, nil
	case model.Move:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Move{
				Move: &proto_messages.Move{MoveX: p.MoveX(), MoveY: p.MoveY()},
			},
		}, nil
	case model.Click:
		return &proto_messages.Message{
			Msg: &proto_messages.Message_Click{
				Click: &proto_messages.Click{MouseButton: proto_messages.MouseButton(p.MouseButton())},
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown processable type: %T", p)
	}
}
