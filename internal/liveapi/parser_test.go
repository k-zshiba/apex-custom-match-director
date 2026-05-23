package liveapi

import (
	"testing"

	"github.com/koshi/apex-custom-match-director/internal/domain"
)

func TestParseEventKillPayload(t *testing.T) {
	event, err := ParseEvent([]byte(`{
		"event": "player_kill",
		"data": {
			"attacker": {"id": "a", "name": "Alpha"},
			"victim": {"id": "b", "name": "Beta"}
		}
	}`))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}
	if event.Kind != domain.EventPlayerKill {
		t.Fatalf("expected kill event, got %s", event.Kind)
	}
	if event.PlayerID != "a" || event.VictimID != "b" {
		t.Fatalf("unexpected parsed players: %+v", event)
	}
}

func TestParseEventUnknown(t *testing.T) {
	event, err := ParseEvent([]byte(`{"event": "new_future_event"}`))
	if err != nil {
		t.Fatalf("ParseEvent returned error: %v", err)
	}
	if event.Kind != domain.EventUnknown {
		t.Fatalf("expected unknown event, got %s", event.Kind)
	}
}
