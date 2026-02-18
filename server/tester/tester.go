package main

import (
	"easy-rc-server/generated/proto-messages"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func main() {

	addr := "0.0.0.0:50267"
	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	fmt.Println("Waiting for 3 seconds...")
	time.Sleep(3 * time.Second)
	fmt.Println("Finished sleeping, starting the testing!!")

	now1 := time.Now().UnixMilli()
	ping := &proto_messages.Ping{Timestamp: now1}
	message1 := proto_messages.Message{Msg: &proto_messages.Message_Ping{Ping: ping}, Debug: "I AM A PING"}
	data, err := proto.Marshal(&message1)

	if err == nil {
		log.Printf("Writing PING message at: %d\n", now1)
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	now2 := time.Now().UnixMilli()
	click := &proto_messages.Click{Timestamp: now2, MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_LEFT}
	message2 := proto_messages.Message{Msg: &proto_messages.Message_Click{Click: click}, Debug: "BLAH blah foobar"}
	data, err = proto.Marshal(&message2)

	if err == nil {
		log.Printf("Writing CLICK LEFT message at: %d\n", now2)
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	now3 := time.Now().UnixMilli()
	moveUp := &proto_messages.Move{MoveX: 100}
	message3 := proto_messages.Message{Msg: &proto_messages.Message_Move{Move: moveUp}, Debug: "Moving up"}
	data, err = proto.Marshal(&message3)
	if err == nil {
		log.Printf("Writing MOVE UP message at: %d\n", now3)
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	fmt.Println("Waiting for 1 second...")
	time.Sleep(1 * time.Second)

	now4 := time.Now().UnixMilli()
	moveRight := &proto_messages.Move{MoveY: 100}
	message4 := proto_messages.Message{Msg: &proto_messages.Message_Move{Move: moveRight}, Debug: "Moving right"}
	data, err = proto.Marshal(&message4)
	if err == nil {
		log.Printf("Writing MOVE RIGHT message at: %d\n", now4)
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	fmt.Println("Waiting for 1 second...")
	time.Sleep(1 * time.Second)

	now5 := time.Now().UnixMilli()
	moveDown := &proto_messages.Move{MoveX: -100}
	message5 := proto_messages.Message{Msg: &proto_messages.Message_Move{Move: moveDown}, Debug: "Moving down"}
	data, err = proto.Marshal(&message5)
	if err == nil {
		log.Printf("Writing MOVE DOWN message at: %d\n", now5)
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	fmt.Println("Waiting for 1 second...")
	time.Sleep(1 * time.Second)

	now6 := time.Now().UnixMilli()
	rightClick := &proto_messages.Click{MouseButton: proto_messages.MouseButton_MOUSE_BUTTON_RIGHT, Timestamp: now6}
	message6 := proto_messages.Message{Msg: &proto_messages.Message_Click{Click: rightClick}, Debug: "Right clicking"}
	data, err = proto.Marshal(&message6)
	if err == nil {
		log.Printf("Writing RIGHT CLICK message at: %d\n", now5)
		c.WriteMessage(websocket.BinaryMessage, data)
	}
}
