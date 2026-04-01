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

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		response, err := handler.HandleMessage(message)
		if err != nil {
			log.Println("Error handling message:", err)
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
