package main

import (
	"log"
	"mouse-server/server/generated"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {

	addr := "192.168.20.16:51022"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	//now1 := time.Now().UnixMilli()
	//ping := &generated.Ping{Timestamp: now1}
	//message1 := generated.Message{Msg: &generated.Message_Ping{Ping: ping}, Debug: "I AM A PING"}
	//data, err := proto.Marshal(&message1)
	//
	//if err == nil {
	//	log.Printf("Writing PING message at: %d\n", now1)
	//	c.WriteMessage(websocket.BinaryMessage, data)
	//}

	now2 := time.Now().UnixMilli()
	click := &generated.Click{Timestamp: now2, MouseButton: generated.MouseButton_MOUSE_BUTTON_LEFT}
	message2 := generated.Message{Msg: &generated.Message_Click{Click: click}, Debug: "BLAH blah foobar"}
	data, err := proto.Marshal(&message2)

	if err == nil {
		log.Printf("Writing CLICK LEFT message at: %d\n", now2)
		c.WriteMessage(websocket.BinaryMessage, data)
	}
}
