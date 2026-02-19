package websocket

import (
	"easy-rc-server/actions"
	"log"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var (
	localNetworks = []net.IPNet{
		{IP: net.IP{10, 0, 0, 0}, Mask: net.CIDRMask(8, 32)},
		{IP: net.IP{172, 16, 0, 0}, Mask: net.CIDRMask(12, 32)},
		{IP: net.IP{192, 168, 0, 0}, Mask: net.CIDRMask(16, 32)},
		{IP: net.IP{127, 0, 0, 0}, Mask: net.CIDRMask(8, 32)},
	}

	upgrader = websocket.Upgrader{CheckOrigin: checkLocalOrigin}
)

func checkLocalOrigin(r *http.Request) bool {
	origin := r.Header.Get("Origin")
	if origin == "" {
		return true
	}

	u, err := url.Parse(origin)
	if err != nil {
		return false
	}

	host := u.Hostname()
	if host == "localhost" {
		return true
	}

	ip := net.ParseIP(host)
	if ip == nil {
		return false
	}

	for _, network := range localNetworks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

func Server(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	handler := NewMessageHandler(actions.NewRobotGo())
	ch := make(chan actions.Processable, 64)

	go func() {
		defer close(ch)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			parsed, err := handler.ParseMessage(message)
			if err != nil {
				log.Println("Error parsing message:", err)
				continue
			}
			ch <- parsed
		}
	}()

	for first := range ch {
		batch := drainChannel(first, ch)
		coalesced := CoalesceActions(batch)

		for _, action := range coalesced {
			response, err := handler.ProcessAction(action)
			if err != nil {
				log.Println("Error processing action:", err)
				continue
			}
			if response != nil {
				log.Println("Sending response on socket:", response)
				if err := c.WriteMessage(websocket.BinaryMessage, response); err != nil {
					log.Println("Error writing response:", err)
					return
				}
			}
		}
	}
}
