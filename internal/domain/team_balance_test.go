package domain

import "testing"

func TestBalanceTeamsMinimizesAverageSpreadForSmallSet(t *testing.T) {
	rankings := []Ranking{
		{Rank: 1, PlayerID: "p1", PlayerName: "P1", Score: 100},
		{Rank: 2, PlayerID: "p2", PlayerName: "P2", Score: 90},
		{Rank: 3, PlayerID: "p3", PlayerName: "P3", Score: 20},
		{Rank: 4, PlayerID: "p4", PlayerName: "P4", Score: 10},
	}

	teams, err := BalanceTeams(rankings, 2, 2)
	if err != nil {
		t.Fatalf("BalanceTeams returned error: %v", err)
	}

	if len(teams) != 2 {
		t.Fatalf("expected 2 teams, got %d", len(teams))
	}
	if teams[0].AverageScore != 55 || teams[1].AverageScore != 55 {
		t.Fatalf("expected balanced averages of 55, got %.1f and %.1f", teams[0].AverageScore, teams[1].AverageScore)
	}
}

func TestBalanceTeamsSupportsUnevenPlayerCounts(t *testing.T) {
	rankings := []Ranking{
		{Rank: 1, PlayerID: "p1", PlayerName: "P1", Score: 100},
		{Rank: 2, PlayerID: "p2", PlayerName: "P2", Score: 80},
		{Rank: 3, PlayerID: "p3", PlayerName: "P3", Score: 60},
		{Rank: 4, PlayerID: "p4", PlayerName: "P4", Score: 40},
		{Rank: 5, PlayerID: "p5", PlayerName: "P5", Score: 20},
	}

	teams, err := BalanceTeams(rankings, 3, 2)
	if err != nil {
		t.Fatalf("BalanceTeams returned error: %v", err)
	}

	totalPlayers := 0
	for _, team := range teams {
		if len(team.Players) > 2 {
			t.Fatalf("team %s exceeded size cap", team.Name)
		}
		totalPlayers += len(team.Players)
	}
	if totalPlayers != len(rankings) {
		t.Fatalf("expected all players assigned, got %d", totalPlayers)
	}
}

func TestBalanceTeamsDeterministic(t *testing.T) {
	rankings := []Ranking{
		{Rank: 1, PlayerID: "a", PlayerName: "A", Score: 100},
		{Rank: 2, PlayerID: "b", PlayerName: "B", Score: 90},
		{Rank: 3, PlayerID: "c", PlayerName: "C", Score: 80},
		{Rank: 4, PlayerID: "d", PlayerName: "D", Score: 70},
		{Rank: 5, PlayerID: "e", PlayerName: "E", Score: 60},
		{Rank: 6, PlayerID: "f", PlayerName: "F", Score: 50},
	}

	first, err := BalanceTeams(rankings, 3, 2)
	if err != nil {
		t.Fatalf("BalanceTeams returned error: %v", err)
	}
	second, err := BalanceTeams(rankings, 3, 2)
	if err != nil {
		t.Fatalf("BalanceTeams returned error: %v", err)
	}

	for i := range first {
		for j := range first[i].Players {
			if first[i].Players[j].PlayerID != second[i].Players[j].PlayerID {
				t.Fatalf("expected deterministic output, got %+v and %+v", first, second)
			}
		}
	}
}
