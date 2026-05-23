package discord

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/koshi/apex-custom-match-director/internal/domain"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (c *Client) PostResult(webhookURL string, state domain.MatchState, teams []domain.Team) error {
	if strings.TrimSpace(webhookURL) == "" {
		return nil
	}
	payload := map[string]string{"content": FormatResult(state, teams)}
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("encode discord payload: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, webhookURL, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("create discord request: %w", RedactSecrets(err))
	}
	req.Header.Set("Content-Type", "application/json")
	res, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("post discord webhook: %w", RedactSecrets(err))
	}
	defer res.Body.Close()
	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return fmt.Errorf("post discord webhook: status %d", res.StatusCode)
	}
	return nil
}

func FormatResult(state domain.MatchState, teams []domain.Team) string {
	var b strings.Builder
	b.WriteString("**Custom Match Results**\n")
	for _, ranking := range state.Rankings {
		fmt.Fprintf(&b, "%d. %s - %d kills\n", ranking.Rank, ranking.PlayerName, ranking.Kills)
	}
	if len(teams) > 0 {
		b.WriteString("\n**Next Teams**\n")
		for _, team := range teams {
			fmt.Fprintf(&b, "%s", team.Name)
			for i, player := range team.Players {
				if i == 0 {
					b.WriteString(": ")
				} else {
					b.WriteString(", ")
				}
				b.WriteString(player.PlayerName)
			}
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func RedactSecrets(err error) error {
	if err == nil {
		return nil
	}
	return redactedError{message: Redact(err.Error())}
}

func Redact(value string) string {
	parts := strings.Split(value, "/api/webhooks/")
	if len(parts) < 2 {
		return value
	}
	return parts[0] + "/api/webhooks/[redacted]"
}

type redactedError struct {
	message string
}

func (e redactedError) Error() string {
	return e.message
}
