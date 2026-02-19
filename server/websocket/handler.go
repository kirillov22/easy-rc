package websocket

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"

	"google.golang.org/protobuf/proto"
)

type MessageHandler struct {
	robot actions.Robot
}

func NewMessageHandler(robot actions.Robot) *MessageHandler {
	return &MessageHandler{robot: robot}
}

func (h *MessageHandler) HandleMessage(data []byte) ([]byte, error) {
	m := &proto_messages.Message{}
	if err := proto.Unmarshal(data, m); err != nil {
		return nil, err
	}

	p, err := actions.FromProto(m)
	if err != nil {
		return nil, err
	}

	r, err := actions.Process(p, h.robot)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	resp, ok := r.(actions.Processable)
	if !ok {
		return nil, nil
	}

	protoResp, err := actions.ToProto(resp)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(protoResp)
}
