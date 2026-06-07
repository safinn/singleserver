package singleserver

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestVerifyWebhookSignature(t *testing.T) {
	body := []byte(`{"ok":true}`)
	secret := "top-secret"
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	signature := fmt.Sprintf("sha256=%x", mac.Sum(nil))

	if !VerifyWebhookSignature(secret, body, signature) {
		t.Fatal("expected signature to verify")
	}
	if VerifyWebhookSignature(secret, body, "sha256=bad") {
		t.Fatal("bad signature verified")
	}
}
