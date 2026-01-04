package main

import (
	"easy-rc-server/websocket"
	"flag"
	"fmt"
	"net"
	//"github.com/go-vgo/robotgo"
	"log"
	"net/http"
	//"strconv"
	//"strings"
)

func main() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/ws", websocket.Server)
	//http.HandleFunc("/", home)
	var serverAddress, port = generateAddress()
	var outboundIp = getOutboundIP()

	log.Printf("Starting websocket server at: %s. Outbound address to connect to: %s:%d\n", serverAddress, outboundIp, port)
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

func generateAddress() (string, int) {
	var port, err = getFreePort()
	if err != nil {
		log.Fatal("Failed to get a free port", err)
	}
	return fmt.Sprintf("0.0.0.0:%d", port), port
}

func getFreePort() (port int, err error) {
	var a *net.TCPAddr
	if a, err = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port, nil
		}
	}
	return
}

func getOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
