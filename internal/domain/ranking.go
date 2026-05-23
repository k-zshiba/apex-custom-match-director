package domain

import "sort"

func RankPlayers(players []PlayerStats) []Ranking {
	rankings := make([]Ranking, 0, len(players))
	for _, p := range players {
		rankings = append(rankings, Ranking{
			PlayerID:    p.ID,
			PlayerName:  p.Name,
			Kills:       p.Kills,
			Deaths:      p.Deaths,
			Damage:      p.Damage,
			WeaponLevel: p.WeaponLevel,
			Score:       PlayerScore(p),
		})
	}

	sort.SliceStable(rankings, func(i, j int) bool {
		a := rankings[i]
		b := rankings[j]
		if a.Kills != b.Kills {
			return a.Kills > b.Kills
		}
		if a.WeaponLevel != b.WeaponLevel {
			return a.WeaponLevel > b.WeaponLevel
		}
		if a.Damage != b.Damage {
			return a.Damage > b.Damage
		}
		if a.Deaths != b.Deaths {
			return a.Deaths < b.Deaths
		}
		if a.PlayerName != b.PlayerName {
			return a.PlayerName < b.PlayerName
		}
		return a.PlayerID < b.PlayerID
	})

	for i := range rankings {
		rankings[i].Rank = i + 1
	}
	return rankings
}

func PlayerScore(p PlayerStats) int {
	return p.Kills*1000 + p.WeaponLevel*100 + p.Damage - p.Deaths*25
}
