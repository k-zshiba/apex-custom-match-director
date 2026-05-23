package liveapi

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/koshi/apex-custom-match-director/internal/domain"
)

func ParseEvent(body []byte) (domain.LiveEvent, error) {
	var raw map[string]any
	if err := json.Unmarshal(body, &raw); err != nil {
		return domain.LiveEvent{}, fmt.Errorf("parse liveapi event: %w", err)
	}

	rawType := firstString(raw, "event", "type", "eventName", "name")
	kind := classify(rawType)
	payload := objectValue(raw, "data", "payload")
	if payload == nil {
		payload = raw
	}

	event := domain.LiveEvent{
		Kind:        kind,
		RawType:     rawType,
		MatchID:     firstString(raw, "matchId", "matchID", "match_id", "gameId", "gameID"),
		PlayerID:    firstString(payload, "playerId", "playerID", "player_id", "uid", "nucleusHash"),
		PlayerName:  firstString(payload, "playerName", "player_name", "name", "displayName"),
		VictimID:    firstString(payload, "victimId", "victimID", "victim_id"),
		VictimName:  firstString(payload, "victimName", "victim_name"),
		Kills:       firstInt(payload, "kills", "killCount"),
		Damage:      firstInt(payload, "damage", "damageDealt"),
		WeaponLevel: firstInt(payload, "weaponLevel", "weapon_level", "weaponIndex"),
	}

	mergePlayerObject(payload, &event, "player", false)
	mergePlayerObject(payload, &event, "attacker", false)
	mergePlayerObject(payload, &event, "killer", false)
	mergePlayerObject(payload, &event, "victim", true)
	return event, nil
}

func classify(rawType string) domain.EventKind {
	normalized := strings.ToLower(strings.ReplaceAll(rawType, "_", ""))
	switch {
	case strings.Contains(normalized, "matchstart") || strings.Contains(normalized, "gamestart"):
		return domain.EventMatchStart
	case strings.Contains(normalized, "matchend") || strings.Contains(normalized, "gameend") || strings.Contains(normalized, "winner"):
		return domain.EventMatchEnd
	case strings.Contains(normalized, "kill"):
		return domain.EventPlayerKill
	case strings.Contains(normalized, "damage"):
		return domain.EventPlayerDamage
	case strings.Contains(normalized, "weapon"):
		return domain.EventWeaponAdvance
	default:
		return domain.EventUnknown
	}
}

func mergePlayerObject(payload map[string]any, event *domain.LiveEvent, key string, victim bool) {
	player := objectValue(payload, key)
	if player == nil {
		return
	}
	id := firstString(player, "id", "playerId", "playerID", "uid", "nucleusHash")
	name := firstString(player, "name", "playerName", "displayName")
	if victim {
		if event.VictimID == "" {
			event.VictimID = id
		}
		if event.VictimName == "" {
			event.VictimName = name
		}
		return
	}
	if event.PlayerID == "" {
		event.PlayerID = id
	}
	if event.PlayerName == "" {
		event.PlayerName = name
	}
}

func objectValue(raw map[string]any, keys ...string) map[string]any {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		object, ok := value.(map[string]any)
		if ok {
			return object
		}
	}
	return nil
}

func firstString(raw map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case string:
			if typed != "" {
				return typed
			}
		case float64:
			return fmt.Sprintf("%.0f", typed)
		}
	}
	return ""
}

func firstInt(raw map[string]any, keys ...string) int {
	for _, key := range keys {
		value, ok := raw[key]
		if !ok {
			continue
		}
		switch typed := value.(type) {
		case float64:
			return int(typed)
		case int:
			return typed
		case string:
			var out int
			if _, err := fmt.Sscanf(typed, "%d", &out); err == nil {
				return out
			}
		}
	}
	return 0
}
