package websocket

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"
	"easy-rc-server/protocol"

	"google.golang.org/protobuf/proto"
)

type MessageHandler struct {
	robot actions.Robot
}

func NewMessageHandler(robot actions.Robot) *MessageHandler {
	return &MessageHandler{robot: robot}
}

func (h *MessageHandler) ParseMessage(data []byte) (actions.Processable, error) {
	m := &proto_messages.Message{}
	if err := proto.Unmarshal(data, m); err != nil {
		return nil, err
	}
	return protocol.FromProto(m)
}

func (h *MessageHandler) ProcessAction(p actions.Processable) ([]byte, error) {
	r, err := actions.Process(p, h.robot)
	if err != nil {
		return nil, err
	}

	if r == nil {
		return nil, nil
	}

	protoResp, err := protocol.ToProto(r)
	if err != nil {
		return nil, err
	}

	return proto.Marshal(protoResp)
}

func (h *MessageHandler) HandleMessage(data []byte) ([]byte, error) {
	p, err := h.ParseMessage(data)
	if err != nil {
		return nil, err
	}
	return h.ProcessAction(p)
}
