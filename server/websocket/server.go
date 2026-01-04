package websocket

import (
	"fmt"
	"log"
	"mouse-server/server/generated"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
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

		m := &generated.Message{}
		err = proto.Unmarshal(message, m)

		if err != nil {
			log.Println("Error unmarshalling message", err)
		}

		fmt.Printf("Message from proto is: %s\n", m)

		c := m.GetClick()
		if c != nil {
			fmt.Printf("Received a click. Button was: %s", c.MouseButton)
		}

		ping := m.GetPing()
		if ping != nil {
			fmt.Printf("Received a PING")
		}

		//packet := messages.FromByteArray(message)
		//mappedMsg, err := messages.MapMessage(&packet)
		//if err != nil {
		//	fmt.Println("Error mapping packet to a message!")
		//	return
		//}
		//
		//fmt.Println("Successfully mapped message!", mappedMsg)
		//mappedMsg.Process()
	}
}
