package websocket

import (
	"net/http"
	"testing"
)

func TestCheckLocalOrigin(t *testing.T) {
	tests := []struct {
		name   string
		origin string
		want   bool
	}{
		{"empty origin", "", true},
		{"localhost", "http://localhost:8080", true},
		{"192.168.x.x", "http://192.168.1.100:3000", true},
		{"10.x.x.x", "http://10.0.0.5", true},
		{"172.16.x.x", "http://172.16.0.1", true},
		{"172.31.x.x upper bound", "http://172.31.255.255", true},
		{"172.32.x.x outside range", "http://172.32.0.1", false},
		{"127.0.0.1 loopback", "http://127.0.0.1", true},
		{"public IP", "http://8.8.8.8", false},
		{"non-local hostname", "http://example.com", false},
		{"invalid URL", "://bad", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &http.Request{Header: http.Header{}}
			if tt.origin != "" {
				r.Header.Set("Origin", tt.origin)
			}

			got := checkLocalOrigin(r)
			if got != tt.want {
				t.Errorf("checkLocalOrigin(%q) = %v, want %v", tt.origin, got, tt.want)
			}
		})
	}
}
