package domain

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

func BalanceTeams(rankings []Ranking, teamCount int, teamSizeCap int) ([]Team, error) {
	if teamCount < 2 || teamCount > 4 {
		return nil, errors.New("team count must be between 2 and 4")
	}
	if teamSizeCap <= 0 {
		return nil, errors.New("team size cap must be greater than zero")
	}
	if len(rankings) == 0 {
		return emptyTeams(teamCount), nil
	}
	if len(rankings) > teamCount*teamSizeCap {
		return nil, fmt.Errorf("%d players exceed capacity %d", len(rankings), teamCount*teamSizeCap)
	}

	players := append([]Ranking(nil), rankings...)
	sort.SliceStable(players, func(i, j int) bool {
		if players[i].Score != players[j].Score {
			return players[i].Score > players[j].Score
		}
		if players[i].Rank != players[j].Rank {
			return players[i].Rank < players[j].Rank
		}
		return players[i].PlayerID < players[j].PlayerID
	})

	targetSizes := targetTeamSizes(len(players), teamCount, teamSizeCap)
	if len(players) > 12 {
		return buildTeams(greedyBalance(players, targetSizes)), nil
	}

	best := balanceSearch{
		players:     players,
		targetSizes: targetSizes,
		teams:       make([][]Ranking, teamCount),
		bestSpread:  math.Inf(1),
		bestPenalty: math.MaxInt,
	}
	best.assign(0)
	return buildTeams(best.bestTeams), nil
}

type balanceSearch struct {
	players     []Ranking
	targetSizes []int
	teams       [][]Ranking
	bestTeams   [][]Ranking
	bestSpread  float64
	bestPenalty int
}

func (s *balanceSearch) assign(index int) {
	if index == len(s.players) {
		if !sizesMatch(s.teams, s.targetSizes) {
			return
		}
		spread := averageSpread(s.teams)
		penalty := deterministicPenalty(s.teams)
		if spread < s.bestSpread || (spread == s.bestSpread && penalty < s.bestPenalty) {
			s.bestSpread = spread
			s.bestPenalty = penalty
			s.bestTeams = cloneTeams(s.teams)
		}
		return
	}

	for teamIndex := range s.teams {
		if len(s.teams[teamIndex]) >= s.targetSizes[teamIndex] {
			continue
		}
		s.teams[teamIndex] = append(s.teams[teamIndex], s.players[index])
		s.assign(index + 1)
		s.teams[teamIndex] = s.teams[teamIndex][:len(s.teams[teamIndex])-1]
	}
}

func targetTeamSizes(playerCount int, teamCount int, cap int) []int {
	sizes := make([]int, teamCount)
	for i := 0; i < playerCount; i++ {
		sizes[i%teamCount]++
	}
	for i := range sizes {
		if sizes[i] > cap {
			sizes[i] = cap
		}
	}
	sort.SliceStable(sizes, func(i, j int) bool { return sizes[i] > sizes[j] })
	return sizes
}

func sizesMatch(teams [][]Ranking, target []int) bool {
	for i := range teams {
		if len(teams[i]) != target[i] {
			return false
		}
	}
	return true
}

func averageSpread(teams [][]Ranking) float64 {
	minAvg := math.Inf(1)
	maxAvg := math.Inf(-1)
	for _, team := range teams {
		if len(team) == 0 {
			continue
		}
		total := 0
		for _, p := range team {
			total += p.Score
		}
		avg := float64(total) / float64(len(team))
		minAvg = math.Min(minAvg, avg)
		maxAvg = math.Max(maxAvg, avg)
	}
	if math.IsInf(minAvg, 0) {
		return 0
	}
	return maxAvg - minAvg
}

func deterministicPenalty(teams [][]Ranking) int {
	penalty := 0
	for teamIndex, team := range teams {
		for playerIndex, p := range team {
			penalty += (teamIndex + 1) * (playerIndex + 1) * p.Rank
		}
	}
	return penalty
}

func cloneTeams(teams [][]Ranking) [][]Ranking {
	out := make([][]Ranking, len(teams))
	for i := range teams {
		out[i] = append([]Ranking(nil), teams[i]...)
	}
	return out
}

func buildTeams(groups [][]Ranking) []Team {
	teams := make([]Team, len(groups))
	for i, group := range groups {
		total := 0
		for _, p := range group {
			total += p.Score
		}
		avg := 0.0
		if len(group) > 0 {
			avg = float64(total) / float64(len(group))
		}
		teams[i] = Team{
			Name:         fmt.Sprintf("Team %d", i+1),
			Players:      append([]Ranking{}, group...),
			TotalScore:   total,
			AverageScore: avg,
		}
	}
	return teams
}

func emptyTeams(teamCount int) []Team {
	teams := make([]Team, teamCount)
	for i := range teams {
		teams[i] = Team{Name: fmt.Sprintf("Team %d", i+1), Players: []Ranking{}}
	}
	return teams
}

func greedyBalance(players []Ranking, targetSizes []int) [][]Ranking {
	teams := make([][]Ranking, len(targetSizes))
	totals := make([]int, len(targetSizes))
	for _, player := range players {
		bestTeam := -1
		bestProjected := math.Inf(1)
		for i := range teams {
			if len(teams[i]) >= targetSizes[i] {
				continue
			}
			projected := projectedSpread(totals, teams, targetSizes, i, player.Score)
			if projected < bestProjected || (projected == bestProjected && i < bestTeam) || bestTeam == -1 {
				bestTeam = i
				bestProjected = projected
			}
		}
		teams[bestTeam] = append(teams[bestTeam], player)
		totals[bestTeam] += player.Score
	}

	improved := true
	for improved {
		improved = false
		current := averageSpread(teams)
		for i := range teams {
			for j := i + 1; j < len(teams); j++ {
				for a := range teams[i] {
					for b := range teams[j] {
						teams[i][a], teams[j][b] = teams[j][b], teams[i][a]
						next := averageSpread(teams)
						if next < current {
							current = next
							improved = true
							continue
						}
						teams[i][a], teams[j][b] = teams[j][b], teams[i][a]
					}
				}
			}
		}
	}
	return teams
}

func projectedSpread(totals []int, teams [][]Ranking, targetSizes []int, teamIndex int, score int) float64 {
	minAvg := math.Inf(1)
	maxAvg := math.Inf(-1)
	for i := range teams {
		total := totals[i]
		count := len(teams[i])
		if i == teamIndex {
			total += score
			count++
		}
		if count == 0 {
			continue
		}
		avg := float64(total) / float64(count)
		minAvg = math.Min(minAvg, avg)
		maxAvg = math.Max(maxAvg, avg)
	}
	if math.IsInf(minAvg, 0) {
		return 0
	}
	return maxAvg - minAvg
}
