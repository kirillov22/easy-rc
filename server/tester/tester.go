//go:build debug

package main

import (
	"easy-rc-server/generated/proto-messages"
	"flag"
	"log"
	"math"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

func circleDelta(t, dt, r float64) (int32, int32) {
	dx := r * (math.Cos(t+dt) - math.Cos(t))
	dy := r * (math.Sin(t+dt) - math.Sin(t))
	return int32(math.Round(dx)), int32(math.Round(dy))
}

func figure8Delta(t, dt, r float64) (int32, int32) {
	dx := r * (math.Cos(t+dt) - math.Cos(t))
	dy := r * (math.Sin(2*(t+dt)) - math.Sin(2*t))
	return int32(math.Round(dx)), int32(math.Round(dy))
}

func main() {
	pattern := flag.String("pattern", "circle", "Motion pattern: circle or figure8")
	interval := flag.Duration("interval", 5*time.Millisecond, "Interval between move messages")
	duration := flag.Duration("duration", 10*time.Second, "Total duration of the stress test")
	radius := flag.Float64("radius", 50.0, "Radius of the motion pattern")
	addr := flag.String("addr", "0.0.0.0:50392", "Server address")
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	log.Println("Waiting for 3 seconds...")
	time.Sleep(3 * time.Second)
	log.Println("Connection ready, verifying with ping/pong...")

	now := time.Now().UnixMilli()
	ping := &proto_messages.Ping{Timestamp: now}
	pingMsg := proto_messages.Message{Msg: &proto_messages.Message_Ping{Ping: ping}, Debug: "stress-test-ping"}
	data, err := proto.Marshal(&pingMsg)
	if err == nil {
		c.WriteMessage(websocket.BinaryMessage, data)
	}

	_, response, err := c.ReadMessage()
	if err != nil {
		log.Fatal("Failed to read PONG response:", err)
	}
	pongMsg := &proto_messages.Message{}
	if err = proto.Unmarshal(response, pongMsg); err != nil {
		log.Fatal("Failed to unmarshal PONG:", err)
	}
	if _, ok := pongMsg.GetMsg().(*proto_messages.Message_Pong); !ok {
		log.Fatalf("Expected PONG but got: %T", pongMsg.GetMsg())
	}
	log.Println("Ping/pong OK. Starting stress test...")

	log.Printf("Pattern=%s, Interval=%v, Duration=%v, Radius=%.0f", *pattern, *interval, *duration, *radius)

	dt := 0.05
	var t float64
	deadline := time.After(*duration)
	ticker := time.NewTicker(*interval)
	defer ticker.Stop()

	sent := 0
loop:
	for {
		select {
		case <-deadline:
			break loop
		case <-ticker.C:
			var dx, dy int32
			switch *pattern {
			case "figure8":
				dx, dy = figure8Delta(t, dt, *radius)
			default:
				dx, dy = circleDelta(t, dt, *radius)
			}
			t += dt

			if dx == 0 && dy == 0 {
				continue
			}

			msg := &proto_messages.Message{
				Msg: &proto_messages.Message_Move{
					Move: &proto_messages.Move{
						Timestamp: time.Now().UnixMilli(),
						MoveX:     dx,
						MoveY:     dy,
					},
				},
			}
			data, err := proto.Marshal(msg)
			if err != nil {
				log.Printf("marshal error: %v", err)
				continue
			}
			if err := c.WriteMessage(websocket.BinaryMessage, data); err != nil {
				log.Fatalf("write error: %v", err)
			}
			sent++
		}
	}

	log.Printf("Stress test complete. Sent %d move messages over %v", sent, *duration)
}
