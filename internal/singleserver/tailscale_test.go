package singleserver

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestDoctorTailscaleKeyExpiry(t *testing.T) {
	soon := time.Now().Add(60 * 24 * time.Hour).Format(time.RFC3339)
	cases := []struct {
		name string
		self *tailscaleSelf
		want string
	}{
		{"tagged node never expires", &tailscaleSelf{Tags: []string{"tag:server"}, KeyExpiry: soon}, "key expiry\tok"},
		{"expiry disabled (zero time)", &tailscaleSelf{KeyExpiry: "0001-01-01T00:00:00Z"}, "key expiry\tok"},
		{"no expiry field", &tailscaleSelf{}, "key expiry\tok"},
		{"untagged with future expiry warns", &tailscaleSelf{KeyExpiry: soon}, "key expiry\tpending"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			var buf bytes.Buffer
			doctorTailscaleKeyExpiry(&buf, &tailscaleStatus{Self: c.self})
			if !strings.Contains(buf.String(), c.want) {
				t.Fatalf("expected %q, got: %q", c.want, buf.String())
			}
		})
	}
}
