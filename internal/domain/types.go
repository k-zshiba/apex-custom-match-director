package domain

type Settings struct {
	LiveAPIPort       int    `json:"liveApiPort"`
	TeamCount         int    `json:"teamCount"`
	TeamSizeCap       int    `json:"teamSizeCap"`
	DiscordWebhookURL string `json:"discordWebhookUrl"`
	AutoPostDiscord   bool   `json:"autoPostDiscord"`
}

type PublicSettings struct {
	LiveAPIPort       int    `json:"liveApiPort"`
	TeamCount         int    `json:"teamCount"`
	TeamSizeCap       int    `json:"teamSizeCap"`
	DiscordWebhookURL string `json:"discordWebhookUrl"`
	AutoPostDiscord   bool   `json:"autoPostDiscord"`
}

func DefaultSettings() Settings {
	return Settings{
		LiveAPIPort:     7777,
		TeamCount:       4,
		TeamSizeCap:     5,
		AutoPostDiscord: true,
	}
}

func (s Settings) Public() PublicSettings {
	return PublicSettings(s)
}

func SettingsFromPublic(s PublicSettings) Settings {
	return Settings(s)
}

func NormalizeSettings(s Settings) Settings {
	if s.LiveAPIPort <= 0 || s.LiveAPIPort > 65535 {
		s.LiveAPIPort = 7777
	}
	if s.TeamCount < 2 {
		s.TeamCount = 2
	}
	if s.TeamCount > 4 {
		s.TeamCount = 4
	}
	if s.TeamSizeCap <= 0 {
		s.TeamSizeCap = 5
	}
	return s
}

type EventKind string

const (
	EventUnknown       EventKind = "unknown"
	EventMatchStart    EventKind = "match_start"
	EventMatchEnd      EventKind = "match_end"
	EventPlayerKill    EventKind = "player_kill"
	EventPlayerDamage  EventKind = "player_damage"
	EventWeaponAdvance EventKind = "weapon_advance"
)

type LiveEvent struct {
	Kind        EventKind `json:"kind"`
	MatchID     string    `json:"matchId,omitempty"`
	PlayerID    string    `json:"playerId,omitempty"`
	PlayerName  string    `json:"playerName,omitempty"`
	VictimID    string    `json:"victimId,omitempty"`
	VictimName  string    `json:"victimName,omitempty"`
	Kills       int       `json:"kills,omitempty"`
	Damage      int       `json:"damage,omitempty"`
	WeaponLevel int       `json:"weaponLevel,omitempty"`
	RawType     string    `json:"rawType,omitempty"`
}

type PlayerStats struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Kills       int    `json:"kills"`
	Deaths      int    `json:"deaths"`
	Damage      int    `json:"damage"`
	WeaponLevel int    `json:"weaponLevel"`
}

type MatchState struct {
	MatchID       string        `json:"matchId"`
	StartedAt     string        `json:"startedAt,omitempty"`
	CompletedAt   string        `json:"completedAt,omitempty"`
	Players       []PlayerStats `json:"players"`
	Rankings      []Ranking     `json:"rankings"`
	Teams         []Team        `json:"teams"`
	UnknownEvents int           `json:"unknownEvents"`
	Completed     bool          `json:"completed"`
	Posted        bool          `json:"posted"`
	LastError     string        `json:"lastError,omitempty"`
}

type Ranking struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"playerId"`
	PlayerName  string `json:"playerName"`
	Kills       int    `json:"kills"`
	Deaths      int    `json:"deaths"`
	Damage      int    `json:"damage"`
	WeaponLevel int    `json:"weaponLevel"`
	Score       int    `json:"score"`
}

type Team struct {
	Name         string    `json:"name"`
	Players      []Ranking `json:"players"`
	AverageScore float64   `json:"averageScore"`
	TotalScore   int       `json:"totalScore"`
}

type ReceiverStatus struct {
	Running bool   `json:"running"`
	Address string `json:"address"`
}
