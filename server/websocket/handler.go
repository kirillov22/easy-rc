package websocket

import (
	"easy-rc-server/actions"
	proto_messages "easy-rc-server/generated/proto-messages"
	"easy-rc-server/protocol"
	"sync"

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

func (h *MessageHandler) ProcessAction(p actions.Processable, wg *sync.WaitGroup, ch chan<- actions.ActionResponse) {
	defer wg.Done()

	r, err := actions.Process(p, h.robot)
	if err != nil {
		ch <- actions.NewActionResponse(nil, err)
		return
	}

	if r == nil {
		ch <- actions.NoopActionResponse()
		return
	}

	protoResp, err := protocol.ToProto(r)
	if err != nil {
		ch <- actions.NewActionResponse(nil, err)
		return
	}

	data, err := proto.Marshal(protoResp)
	if err != nil {
		ch <- actions.NewActionResponse(nil, err)
		return
	}

	ch <- actions.NewActionResponse(data, nil)
}

//func (h *MessageHandler) HandleMessage(data []byte) actions.ActionResponse {
//	p, err := h.ParseMessage(data)
//	if err != nil {
//		return actions.NewActionResponse(nil, err)
//	}
//	return h.ProcessAction(p)
//}
