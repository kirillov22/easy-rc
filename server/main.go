package main

import (
	ws "easy-rc-server/websocket"
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	flag.Parse()
	http.HandleFunc("/ws", ws.Server)

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	outboundIP := getOutboundIP()
	log.Printf("Starting websocket server at: 0.0.0.0:%d. Outbound address to connect to: %s:%d\n", port, outboundIP, port)

	server := &http.Server{}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down", sig)
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
	}()

	if err := server.Serve(listener); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
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
