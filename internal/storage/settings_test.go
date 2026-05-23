package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/koshi/apex-custom-match-director/internal/domain"
)

func TestSettingsStoreRoundTrip(t *testing.T) {
	store, err := NewSettingsStore(filepath.Join(t.TempDir(), "settings.sqlite"))
	if err != nil {
		t.Fatalf("NewSettingsStore returned error: %v", err)
	}
	defer store.Close()

	want := domain.DefaultSettings()
	want.LiveAPIPort = 8787
	want.TeamCount = 3
	want.DiscordWebhookURL = "https://discord.com/api/webhooks/1/secret"
	if err := store.Save(context.Background(), want); err != nil {
		t.Fatalf("Save returned error: %v", err)
	}

	got, err := store.Load(context.Background())
	if err != nil {
		t.Fatalf("Load returned error: %v", err)
	}
	if got.LiveAPIPort != want.LiveAPIPort || got.TeamCount != want.TeamCount || got.DiscordWebhookURL != want.DiscordWebhookURL {
		t.Fatalf("settings mismatch: got %+v want %+v", got, want)
	}
}
