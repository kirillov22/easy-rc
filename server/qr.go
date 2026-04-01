package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
)

func saveQRImage(url string) string {
	qr, err := qrcode.New(url)
	if err != nil {
		log.Printf("Failed to generate QR code image: %v", err)
		return ""
	}

	path := filepath.Join(os.TempDir(), "easy-rc-qr.png")
	w, err := standard.New(path,
		standard.WithBuiltinImageEncoder(standard.PNG_FORMAT),
		standard.WithQRWidth(10),
		standard.WithBorderWidth(20),
	)
	if err != nil {
		log.Printf("Failed to create QR image writer: %v", err)
		return ""
	}

	if err := qr.Save(w); err != nil {
		log.Printf("Failed to save QR image: %v", err)
		return ""
	}

	return path
}

func removeQRImage(path string) {
	if path == "" {
		return
	}
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to remove QR image: %v", err)
	}
}
