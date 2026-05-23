package domain

import (
	"sort"
	"sync"
	"time"
)

type Aggregator struct {
	mu    sync.Mutex
	state MatchState
}

func NewAggregator() *Aggregator {
	return &Aggregator{}
}

func (a *Aggregator) Apply(event LiveEvent) MatchState {
	a.mu.Lock()
	defer a.mu.Unlock()

	now := time.Now().UTC()
	if event.MatchID != "" {
		a.state.MatchID = event.MatchID
	}

	switch event.Kind {
	case EventMatchStart:
		a.state = MatchState{MatchID: event.MatchID, StartedAt: now.Format(time.RFC3339)}
	case EventMatchEnd:
		a.state.Completed = true
		a.state.CompletedAt = now.Format(time.RFC3339)
	case EventPlayerKill:
		killer := a.ensurePlayer(event.PlayerID, event.PlayerName)
		killer.Kills += max(1, event.Kills)
		if event.WeaponLevel > killer.WeaponLevel {
			killer.WeaponLevel = event.WeaponLevel
		}
		if event.VictimID != "" || event.VictimName != "" {
			victim := a.ensurePlayer(event.VictimID, event.VictimName)
			victim.Deaths++
		}
	case EventPlayerDamage:
		player := a.ensurePlayer(event.PlayerID, event.PlayerName)
		player.Damage += max(0, event.Damage)
	case EventWeaponAdvance:
		player := a.ensurePlayer(event.PlayerID, event.PlayerName)
		if event.WeaponLevel > player.WeaponLevel {
			player.WeaponLevel = event.WeaponLevel
		}
	default:
		a.state.UnknownEvents++
	}

	a.refreshRankingsLocked()
	return a.snapshotLocked()
}

func (a *Aggregator) State() MatchState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.refreshRankingsLocked()
	return a.snapshotLocked()
}

func (a *Aggregator) MarkPosted() MatchState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.Posted = true
	return a.snapshotLocked()
}

func (a *Aggregator) SetTeams(teams []Team) MatchState {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.Teams = append([]Team(nil), teams...)
	return a.snapshotLocked()
}

func (a *Aggregator) SetLastError(message string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.state.LastError = message
}

func (a *Aggregator) ensurePlayer(id string, name string) *PlayerStats {
	if id == "" {
		id = name
	}
	if name == "" {
		name = id
	}
	for i := range a.state.Players {
		if a.state.Players[i].ID == id {
			if a.state.Players[i].Name == "" {
				a.state.Players[i].Name = name
			}
			return &a.state.Players[i]
		}
	}
	a.state.Players = append(a.state.Players, PlayerStats{ID: id, Name: name})
	return &a.state.Players[len(a.state.Players)-1]
}

func (a *Aggregator) refreshRankingsLocked() {
	sort.SliceStable(a.state.Players, func(i, j int) bool {
		if a.state.Players[i].Name != a.state.Players[j].Name {
			return a.state.Players[i].Name < a.state.Players[j].Name
		}
		return a.state.Players[i].ID < a.state.Players[j].ID
	})
	a.state.Rankings = RankPlayers(a.state.Players)
}

func (a *Aggregator) snapshotLocked() MatchState {
	out := a.state
	out.Players = append([]PlayerStats(nil), a.state.Players...)
	out.Rankings = append([]Ranking(nil), a.state.Rankings...)
	out.Teams = append([]Team(nil), a.state.Teams...)
	if out.Players == nil {
		out.Players = []PlayerStats{}
	}
	if out.Rankings == nil {
		out.Rankings = []Ranking{}
	}
	if out.Teams == nil {
		out.Teams = []Team{}
	}
	return out
}
