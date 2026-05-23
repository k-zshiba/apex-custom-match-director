package domain

import "testing"

func TestRankPlayersKillCenteredWithDeterministicTies(t *testing.T) {
	players := []PlayerStats{
		{ID: "b", Name: "Beta", Kills: 4, Deaths: 1, Damage: 900, WeaponLevel: 8},
		{ID: "a", Name: "Alpha", Kills: 5, Deaths: 2, Damage: 700, WeaponLevel: 6},
		{ID: "c", Name: "Charlie", Kills: 5, Deaths: 1, Damage: 600, WeaponLevel: 6},
	}

	got := RankPlayers(players)

	if got[0].PlayerName != "Alpha" {
		t.Fatalf("expected Alpha first by damage tie-breaker, got %s", got[0].PlayerName)
	}
	if got[1].PlayerName != "Charlie" {
		t.Fatalf("expected Charlie second, got %s", got[1].PlayerName)
	}
	if got[2].PlayerName != "Beta" {
		t.Fatalf("expected Beta third, got %s", got[2].PlayerName)
	}
}

func TestAggregatorUnknownEventsDoNotCrash(t *testing.T) {
	agg := NewAggregator()

	state := agg.Apply(LiveEvent{Kind: EventUnknown, RawType: "future_event"})

	if state.UnknownEvents != 1 {
		t.Fatalf("expected unknown event count 1, got %d", state.UnknownEvents)
	}
}

func TestAggregatorAppliesKillsAndDeaths(t *testing.T) {
	agg := NewAggregator()
	state := agg.Apply(LiveEvent{
		Kind:       EventPlayerKill,
		PlayerID:   "killer",
		PlayerName: "Killer",
		VictimID:   "victim",
		VictimName: "Victim",
	})

	if len(state.Rankings) != 2 {
		t.Fatalf("expected two players, got %d", len(state.Rankings))
	}
	if state.Rankings[0].PlayerName != "Killer" || state.Rankings[0].Kills != 1 {
		t.Fatalf("expected killer ranked first with one kill, got %+v", state.Rankings[0])
	}
}
