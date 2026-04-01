package main

import (
	"context"
	"easy-rc-server/actions"
	ws "easy-rc-server/websocket"
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/getlantern/systray"
)

//go:embed static/*
var staticFiles embed.FS

func main() {
	flag.Parse()
	http.Handle("/ws", ws.NewServer(actions.NewRobotGo()))

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

	qrImagePath := saveQRImage(connectURL)

	server := &http.Server{}

	shutdown := func() {
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("Shutdown error: %v", err)
		}
		removeQRImage(qrImagePath)
	}

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigCh
		log.Printf("Received signal %v, shutting down", sig)
		shutdown()
		os.Exit(0)
	}()

	go func() {
		if err := server.Serve(listener); err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	systray.Run(func() {
		systray.SetTitle("RC")
		systray.SetTooltip("Easy-RC Server")

		mShowQR := systray.AddMenuItem("Show QR Code", "Open QR code to scan")
		systray.AddSeparator()
		mQuit := systray.AddMenuItem("Quit", "Shut down the server")

		go func() {
			for {
				select {
				case <-mShowQR.ClickedCh:
					if qrImagePath != "" {
						exec.Command("open", qrImagePath).Start()
					}
				case <-mQuit.ClickedCh:
					shutdown()
					systray.Quit()
				}
			}
		}()
	}, nil)
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
