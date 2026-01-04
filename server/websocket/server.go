package websocket

import (
	"easy-rc-server/actions"
	"easy-rc-server/generated/proto-messages"
	//"fmt"
	"log"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"

	//"google.golang.org/protobuf/proto"

	//"mouse-server/messages"
	"net/http"
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
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}

		m := &proto_messages.Message{}
		err = proto.Unmarshal(message, m)

		if err != nil {
			log.Println("Error unmarshalling message", err)
		}

		// TODO Handle error
		p, err := actions.FromProto(m)
		r, err := actions.Process(p)

		if r != nil {
			// TODO: Handle sending back a response
		}
	}
}
