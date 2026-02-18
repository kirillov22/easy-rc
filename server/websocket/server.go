package websocket

import (
	"easy-rc-server/actions"
	"easy-rc-server/generated/proto-messages"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{} // use default options

func Server(w http.ResponseWriter, r *http.Request) {
	// TODO: Change this later
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// TODO: Extract this out into a separate file so it can be tested
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)

		m := &proto_messages.Message{}
		err = proto.Unmarshal(message, m)

		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}

		p, err := actions.FromProto(m)
		if err != nil {
			log.Println("Error mapping message:", err)
			continue
		}

		r, err := actions.Process(p)
		if err != nil {
			log.Println("Error processing message:", err)
			continue
		}

		if r != nil {
			resp, ok := r.(actions.Processable)
			if !ok {
				log.Println("Process result is not a Processable")
				continue
			}
			protoResp, err := actions.ToProto(resp)
			if err != nil {
				log.Println("Error converting response to proto:", err)
				continue
			}
			data, err := proto.Marshal(protoResp)
			if err != nil {
				log.Println("Error marshalling response:", err)
				continue
			}
			log.Println("Sending response on socket:", data)
			err = c.WriteMessage(websocket.BinaryMessage, data)
			if err != nil {
				log.Println("Error writing response:", err)
				break
			}
		}
	}
}
