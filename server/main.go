package main

import (
	ws "easy-rc-server/websocket"
	"context"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	qrcode "github.com/yeqown/go-qrcode/v2"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	flag.Parse()
	http.HandleFunc("/ws", ws.Server)

	staticContent, err := fs.Sub(staticFiles, "static")
	if err != nil {
		log.Fatalf("Failed to create static sub-filesystem: %v", err)
	}
	http.Handle("/", http.FileServer(http.FS(staticContent)))

	listener, err := net.Listen("tcp", "0.0.0.0:0")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	port := listener.Addr().(*net.TCPAddr).Port
	outboundIP := getOutboundIP()
	connectURL := fmt.Sprintf("http://%s:%d", outboundIP, port)
	log.Printf("Starting websocket server at: 0.0.0.0:%d. Outbound address to connect to: %s:%d\n", port, outboundIP, port)
	printQRCode(connectURL)

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

func printQRCode(url string) {
	qr, err := qrcode.New(url)
	if err != nil {
		log.Printf("Failed to generate QR code: %v", err)
		return
	}

	if err := qr.Save(&stdoutWriter{}); err != nil {
		log.Printf("Failed to print QR code: %v", err)
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
