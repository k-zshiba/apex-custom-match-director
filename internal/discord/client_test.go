package discord

import (
	"strings"
	"testing"
)

func TestRedactWebhookURL(t *testing.T) {
	raw := "post https://discord.com/api/webhooks/123/secret failed"
	got := Redact(raw)
	if strings.Contains(got, "123/secret") {
		t.Fatalf("expected webhook secret to be redacted, got %q", got)
	}
	if !strings.Contains(got, "[redacted]") {
		t.Fatalf("expected redaction marker, got %q", got)
	}
}
